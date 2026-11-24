package httpx

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

// latencyBuckets has an edge at 0.2 s, the NFR-1 target, so the share of
// requests under it is exact rather than interpolated. scripts/p95.sh relies on it.
var latencyBuckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.2, 0.5, 1, 2.5, 5}

// Metrics records how long each request takes, measured at the server (NFR-1).
type Metrics struct {
	duration *prometheus.HistogramVec
}

// NewMetrics registers the request histogram with reg.
func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "vernacular_http_request_duration_seconds",
		Help:    "Time to answer a request, by method, route pattern and status class.",
		Buckets: latencyBuckets,
	}, []string{"method", "route", "status"})}
	reg.MustRegister(m.duration)
	return m
}

// observe records one request's duration. Requests are labelled with chi
// patterns such as /v1/entries/{id}, so the label set stays small; paths that
// match no route share one label.
func (m *Metrics) observe(routes chi.Routes, r *http.Request, status int, d time.Duration) {
	class := strconv.Itoa(status/100) + "xx"
	m.duration.WithLabelValues(r.Method, routeOf(routes, r), class).Observe(d.Seconds())
}

// routeOf is the pattern chi routed r to. A request that middleware answered
// before routing (401, 429, a validation error) is matched here instead.
func routeOf(routes chi.Routes, r *http.Request) string {
	if rc := chi.RouteContext(r.Context()); rc != nil && rc.RoutePattern() != "" {
		return rc.RoutePattern()
	}
	rc := chi.NewRouteContext()
	if routes.Match(rc, r.Method, r.URL.Path) {
		return rc.RoutePattern()
	}
	return "unmatched"
}
