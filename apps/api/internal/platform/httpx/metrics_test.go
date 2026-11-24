package httpx_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// TestMetricsLabelRoutes: requests are counted by route pattern, not raw path,
// and unrouted paths share one label, so the label set cannot grow with traffic.
func TestMetricsLabelRoutes(t *testing.T) {
	reg := prometheus.NewRegistry()
	r := httpx.NewRouter(httpx.RouterOptions{
		Log:     slog.New(slog.NewTextHandler(io.Discard, nil)),
		Metrics: httpx.NewMetrics(reg),
	})
	// Like the rate limiter: answers before routing for one caller.
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Header.Get("X-Limited") != "" {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, req)
		})
	})
	r.Get("/v1/entries/{id}", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	for _, c := range []struct{ path, limited string }{
		{"/v1/entries/a", ""}, {"/v1/entries/b", ""}, {"/v1/entries/c", "yes"}, {"/nope/1", ""}, {"/nope/2", ""},
	} {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, c.path, nil)
		req.Header.Set("X-Limited", c.limited)
		r.ServeHTTP(httptest.NewRecorder(), req)
	}

	got := map[string]int{}
	mfs, err := reg.Gather()
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range mfs[0].GetMetric() {
		labels := map[string]string{}
		for _, l := range m.GetLabel() {
			labels[l.GetName()] = l.GetValue()
		}
		got[labels["route"]+" "+labels["status"]] = int(m.GetHistogram().GetSampleCount())
	}
	wantCounts := map[string]int{"/v1/entries/{id} 2xx": 2, "/v1/entries/{id} 4xx": 1, "unmatched 4xx": 2}
	if len(got) != len(wantCounts) {
		t.Fatalf("series = %v, want %v", got, wantCounts)
	}
	for k, v := range wantCounts {
		if got[k] != v {
			t.Errorf("%s = %d, want %d (all: %v)", k, got[k], v, got)
		}
	}
}
