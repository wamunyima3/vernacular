package accounts

import (
	"context"
	"net/http"
	"strings"

	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// Authenticator resolves a session token or an API key to the caller.
type Authenticator interface {
	Authenticate(ctx context.Context, token string) (Principal, error)
	AuthenticateKey(ctx context.Context, secret string) (Principal, error)
}

// Middleware puts the caller in the request context when an API key
// (X-API-Key) or a Bearer session token is sent. Requests with neither pass
// through anonymous; the contract decides which operations need one. A bad
// key or token gets 401.
func Middleware(auth Authenticator, errs httpx.Errors) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				p   Principal
				err error
			)
			if key := r.Header.Get("X-API-Key"); key != "" {
				p, err = auth.AuthenticateKey(r.Context(), key)
			} else if token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer "); ok && token != "" {
				p, err = auth.Authenticate(r.Context(), token)
			} else {
				next.ServeHTTP(w, r)
				return
			}
			if err != nil {
				errs.ResponseError(w, r, err) // ErrNoSession and ErrNoKey become 401
				return
			}
			next.ServeHTTP(w, r.WithContext(WithPrincipal(r.Context(), p)))
		})
	}
}
