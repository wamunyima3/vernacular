package accounts

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// Authenticator resolves a session token to the caller.
type Authenticator interface {
	Authenticate(ctx context.Context, token string) (Principal, error)
}

// Middleware puts the caller in the request context when a Bearer token is
// sent. Requests without one pass through anonymous; a bad token gets 401.
func Middleware(auth Authenticator, errs httpx.Errors) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || token == "" {
				next.ServeHTTP(w, r)
				return
			}
			p, err := auth.Authenticate(r.Context(), token)
			switch {
			case errors.Is(err, ErrNoSession):
				errs.Write(w, http.StatusUnauthorized, gen.Unauthorized, "Session is invalid or expired")
			case err != nil:
				errs.ResponseError(w, r, err)
			default:
				next.ServeHTTP(w, r.WithContext(WithPrincipal(r.Context(), p)))
			}
		})
	}
}
