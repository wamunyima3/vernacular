package server_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/languages"
	"github.com/wamunyima3/vernacular/apps/api/internal/server"
)

type fakePinger struct{ err error }

func (f fakePinger) Ping(context.Context) error { return f.err }

type fakeLanguages struct {
	langs []languages.Language
	err   error
}

func (f fakeLanguages) List(context.Context) ([]languages.Language, error) { return f.langs, f.err }

type fakeAuth map[string]accounts.Principal

func (f fakeAuth) Authenticate(_ context.Context, token string) (accounts.Principal, error) {
	if token == "boom" {
		return accounts.Principal{}, errors.New("db down")
	}
	p, ok := f[token]
	if !ok {
		return accounts.Principal{}, accounts.ErrNoSession
	}
	return p, nil
}

func newServer(deps server.Deps) http.Handler {
	return server.New(server.Options{
		Log:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		DocsBaseURL: "https://docs.example/errors",
	}, deps)
}

func get(t *testing.T, h http.Handler, path string, out any) int {
	t.Helper()
	return do(t, h, httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil), out)
}

func do(t *testing.T, h http.Handler, req *http.Request, out any) int {
	t.Helper()
	path := req.URL.Path
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("%s: Content-Type = %q, want application/json", path, ct)
	}
	if err := json.Unmarshal(rec.Body.Bytes(), out); err != nil {
		t.Fatalf("%s: decode body %q: %v", path, rec.Body.String(), err)
	}
	return rec.Code
}

func TestListLanguages(t *testing.T) {
	bemba := languages.Language{Code: "bem", Name: "Bemba", OwnName: "Ichibemba", Wave: 1, Status: "active"}

	tests := []struct {
		name       string
		deps       server.Deps
		wantStatus int
		wantCodes  []string
		wantError  gen.ErrorCode
	}{
		{
			name:       "returns languages in store order",
			deps:       server.Deps{Languages: fakeLanguages{langs: []languages.Language{bemba}}},
			wantStatus: http.StatusOK,
			wantCodes:  []string{"bem"},
		},
		{
			name:       "empty list is an empty array, not null",
			deps:       server.Deps{Languages: fakeLanguages{langs: nil}},
			wantStatus: http.StatusOK,
			wantCodes:  []string{},
		},
		{
			name:       "store failure becomes a uniform internal error",
			deps:       server.Deps{Languages: fakeLanguages{err: errors.New("db down")}},
			wantStatus: http.StatusInternalServerError,
			wantError:  gen.ErrorCodeInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body struct {
				gen.LanguageList
				gen.Error
			}
			status := get(t, newServer(tt.deps), "/v1/languages", &body)
			if status != tt.wantStatus {
				t.Fatalf("status = %d, want %d", status, tt.wantStatus)
			}
			if tt.wantError != "" {
				if body.Code != tt.wantError || body.DocsUrl != "https://docs.example/errors#"+string(tt.wantError) {
					t.Fatalf("error body = %+v, want code %q with docs link", body.Error, tt.wantError)
				}
				return
			}
			if body.Data == nil {
				t.Fatal("data is null, want an array")
			}
			got := make([]string, len(body.Data))
			for i, l := range body.Data {
				got[i] = l.Code
			}
			if len(got) != len(tt.wantCodes) || (len(got) > 0 && got[0] != tt.wantCodes[0]) {
				t.Fatalf("codes = %v, want %v", got, tt.wantCodes)
			}
		})
	}
}

func TestHealth(t *testing.T) {
	tests := []struct {
		name       string
		pingErr    error
		wantStatus int
	}{
		{"database reachable", nil, http.StatusOK},
		{"database down", errors.New("refused"), http.StatusServiceUnavailable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body map[string]any
			status := get(t, newServer(server.Deps{DB: fakePinger{err: tt.pingErr}}), "/healthz", &body)
			if status != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body %v)", status, tt.wantStatus, body)
			}
		})
	}
}

func TestUnknownRouteUsesErrorBody(t *testing.T) {
	var body gen.Error
	if status := get(t, newServer(server.Deps{}), "/v1/nope", &body); status != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", status)
	}
	if body.Code != gen.ErrorCodeNotFound {
		t.Fatalf("code = %q, want %q", body.Code, gen.ErrorCodeNotFound)
	}
}

func TestGetMe(t *testing.T) {
	auth := fakeAuth{"good": {
		ID: "u1", Name: "Mwila", Email: "mwila@example.org", Role: gen.RoleContributor,
		Approvals: []accounts.Approval{{Language: "bem", Level: gen.ReviewerLevelReviewer}},
	}}
	tests := []struct {
		name       string
		header     string
		wantStatus int
		wantCode   gen.ErrorCode
	}{
		{"no token", "", http.StatusUnauthorized, gen.ErrorCodeUnauthorized},
		{"unknown token", "Bearer nope", http.StatusUnauthorized, gen.ErrorCodeUnauthorized},
		{"not a bearer token", "Basic abc", http.StatusUnauthorized, gen.ErrorCodeUnauthorized},
		{"lookup failure", "Bearer boom", http.StatusInternalServerError, gen.ErrorCodeInternal},
		{"valid session", "Bearer good", http.StatusOK, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/v1/me", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			var body struct {
				gen.Me
				Code gen.ErrorCode `json:"code"`
			}
			status := do(t, newServer(server.Deps{Accounts: auth}), req, &body)
			if status != tt.wantStatus || body.Code != tt.wantCode {
				t.Fatalf("got %d %q, want %d %q", status, body.Code, tt.wantStatus, tt.wantCode)
			}
			if status == http.StatusOK && (body.Id != "u1" || body.Role != gen.RoleContributor || len(body.Approvals) != 1) {
				t.Errorf("me = %+v", body.Me)
			}
		})
	}
}
