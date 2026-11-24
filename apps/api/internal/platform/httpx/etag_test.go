package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

func TestETag(t *testing.T) {
	h := httpx.ETag(func(*http.Request) string { return "private, max-age=300" })(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/missing" {
			http.Error(w, "no", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"word":"amenshi"}`))
	}))
	serve := func(method, path, ifNoneMatch string) *httptest.ResponseRecorder {
		r := httptest.NewRequestWithContext(t.Context(), method, path, nil)
		if ifNoneMatch != "" {
			r.Header.Set("If-None-Match", ifNoneMatch)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		return rec
	}

	first := serve(http.MethodGet, "/", "")
	tag := first.Header().Get("ETag")
	if first.Code != 200 || tag == "" || first.Body.String() != `{"word":"amenshi"}` || first.Header().Get("Cache-Control") != "private, max-age=300" {
		t.Fatalf("first: %d %q %q %q", first.Code, tag, first.Body, first.Header().Get("Cache-Control"))
	}
	tests := []struct {
		name, method, path, ifNoneMatch string
		want                            int
	}{
		{"matching tag", http.MethodGet, "/", tag, http.StatusNotModified},
		{"strong form of the tag, in a list", http.MethodGet, "/", `"x", ` + tag[2:], http.StatusNotModified},
		{"other tag", http.MethodGet, "/", `W/"stale"`, http.StatusOK},
		{"not found keeps its status", http.MethodGet, "/missing", "*", http.StatusNotFound},
		{"POST is untouched", http.MethodPost, "/", tag, http.StatusOK},
	}
	for _, tt := range tests {
		rec := serve(tt.method, tt.path, tt.ifNoneMatch)
		if rec.Code != tt.want || (rec.Code == http.StatusNotModified && rec.Body.Len() != 0) {
			t.Errorf("%s: %d with %d bytes, want %d", tt.name, rec.Code, rec.Body.Len(), tt.want)
		}
	}
}
