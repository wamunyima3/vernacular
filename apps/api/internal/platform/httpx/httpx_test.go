package httpx_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

func TestRouterErrorsAreUniform(t *testing.T) {
	errs := httpx.Errors{DocsBaseURL: "https://docs.example/errors", Log: slog.New(slog.NewTextHandler(io.Discard, nil))}
	r := httpx.NewRouter(httpx.RouterOptions{Log: errs.Log, Errors: errs})
	r.Get("/panic", func(http.ResponseWriter, *http.Request) { panic("boom") })
	r.Get("/ok", func(http.ResponseWriter, *http.Request) {})

	tests := []struct {
		method, path string
		wantStatus   int
		wantCode     gen.ErrorCode
	}{
		{http.MethodGet, "/panic", http.StatusInternalServerError, gen.ErrorCodeInternal},
		{http.MethodGet, "/missing", http.StatusNotFound, gen.ErrorCodeNotFound},
		{http.MethodPost, "/ok", http.StatusMethodNotAllowed, gen.ErrorCodeMethodNotAllowed},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), tt.method, tt.path, nil))

			var body gen.Error
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode %q: %v", rec.Body.String(), err)
			}
			if rec.Code != tt.wantStatus || body.Code != tt.wantCode {
				t.Errorf("got %d %q, want %d %q", rec.Code, body.Code, tt.wantStatus, tt.wantCode)
			}
			if want := "https://docs.example/errors#" + string(tt.wantCode); body.DocsUrl != want {
				t.Errorf("docs_url = %q, want %q", body.DocsUrl, want)
			}
		})
	}
}
