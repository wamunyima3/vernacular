package httpx

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// RouterOptions configures the shared middleware stack.
type RouterOptions struct {
	Log            *slog.Logger
	AllowedOrigins []string
	Errors         Errors
}

// NewRouter returns a chi router with request IDs, logging, panic recovery,
// CORS and compression (NFR-3) applied to every route.
func NewRouter(opts RouterOptions) chi.Router {
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		requestLogger(opts.Log),
		opts.Errors.Recover,
		cors.Handler(cors.Options{
			AllowedOrigins: opts.AllowedOrigins,
			AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
			AllowedHeaders: []string{"Authorization", "Content-Type", "X-API-Key"},
			ExposedHeaders: []string{"ETag", "X-Request-Id", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset", "Retry-After"},
			MaxAge:         300,
		}),
		middleware.Compress(5),
	)
	r.NotFound(opts.Errors.NotFound)
	r.MethodNotAllowed(opts.Errors.MethodNotAllowed)
	return r
}

func requestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()
			next.ServeHTTP(ww, r)
			log.InfoContext(r.Context(), "request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", middleware.GetReqID(r.Context()),
			)
		})
	}
}
