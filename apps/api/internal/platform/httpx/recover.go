package httpx

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
)

// Recover turns a panic into the standard 500 error body and logs the stack.
// It replaces chi's Recoverer, which answers in plain text.
func (e Errors) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			// The server uses this sentinel to abort a response; let it through.
			if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
				panic(rec)
			}
			e.Log.ErrorContext(r.Context(), "panic", "panic", rec, "path", r.URL.Path, "stack", string(debug.Stack()))
			e.Write(w, http.StatusInternalServerError, gen.ErrorCodeInternal, "Internal server error")
		}()
		next.ServeHTTP(w, r)
	})
}
