package server

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// RateLimits are requests per minute (DEV-3, NFR-5). Zero turns a limit off.
type RateLimits struct {
	PerKey int // per API key, or per signed-in user
	PerIP  int // per client address, for anonymous calls
	IP     httpx.ClientIP
}

// rateLimit counts each call against its API key, its signed-in user, or,
// for anonymous calls, its client address, in one-minute windows. Every
// counted response carries X-RateLimit-Limit, -Remaining and -Reset; going
// over answers 429 with the standard error body and Retry-After (US-13).
// Counters are in memory, which suits the single API process (ADR-0002).
func rateLimit(limits RateLimits, errs httpx.Errors) func(http.Handler) http.Handler {
	tooMany := httprate.WithLimitHandler(func(w http.ResponseWriter, _ *http.Request) {
		errs.Write(w, http.StatusTooManyRequests, gen.ErrorCodeRateLimited, "Too many requests. Wait until the time in X-RateLimit-Reset.")
	})
	limiter := func(n int) *httprate.RateLimiter {
		if n <= 0 {
			return nil
		}
		return httprate.NewRateLimiter(n, time.Minute, tooMany)
	}
	perKey, perIP := limiter(limits.PerKey), limiter(limits.PerIP)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthz" { // probes are never limited
				next.ServeHTTP(w, r)
				return
			}
			var (
				l   *httprate.RateLimiter
				key string
			)
			switch p, ok := accounts.FromContext(r.Context()); {
			case !ok:
				l = perIP
				if l != nil {
					key = "ip:" + limits.IP.Of(r).String()
				}
			case p.ByKey():
				l, key = perKey, "key:"+p.KeyID.String()
			default:
				l, key = perKey, "user:"+p.ID
			}
			if l == nil || !l.RespondOnLimit(w, r, key) {
				next.ServeHTTP(w, r)
			}
		})
	}
}
