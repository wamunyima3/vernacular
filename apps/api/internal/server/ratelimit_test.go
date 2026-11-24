package server_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/server"
)

// TestRateLimit: two anonymous calls a minute per IP; a key has its own count.
func TestRateLimit(t *testing.T) {
	h := server.New(server.Options{
		Log:        slog.New(slog.NewTextHandler(io.Discard, nil)),
		RateLimits: server.RateLimits{PerKey: 5, PerIP: 2},
	}, server.Deps{Accounts: fakeAuth{}})
	call := func(key string) *httptest.ResponseRecorder {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/v1/me", nil)
		if key != "" {
			req.Header.Set("X-API-Key", key)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec
	}

	for want := 1; want >= 0; want-- {
		rec := call("")
		if got := rec.Header().Get("X-RateLimit-Remaining"); got != strconv.Itoa(want) || rec.Header().Get("X-RateLimit-Limit") != "2" {
			t.Fatalf("anonymous call: remaining %q, limit %q; want %d of 2", got, rec.Header().Get("X-RateLimit-Limit"), want)
		}
	}
	rec := call("")
	var body gen.Error
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || rec.Code != http.StatusTooManyRequests || body.Code != gen.ErrorCodeRateLimited {
		t.Fatalf("third call: %d %s, want 429 rate_limited", rec.Code, rec.Body)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Error("429 without Retry-After")
	}
	if rec := call("vk_good"); rec.Code == http.StatusTooManyRequests || rec.Header().Get("X-RateLimit-Limit") != "5" {
		t.Errorf("a key after its IP ran out: %d, limit %q; want its own count of 5", rec.Code, rec.Header().Get("X-RateLimit-Limit"))
	}
}
