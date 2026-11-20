// Package httpx holds HTTP concerns shared by every module: the uniform error
// body and the router middleware stack.
package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
)

// Stable error codes. Clients match on these, so never rename one.
const (
	CodeBadRequest       = "bad_request"
	CodeNotFound         = "not_found"
	CodeMethodNotAllowed = "method_not_allowed"
	CodeInternal         = "internal"
	CodeUnavailable      = "unavailable"
)

// Errors builds the uniform error body {code, message, docs_url} used by every response.
type Errors struct {
	DocsBaseURL string
	Log         *slog.Logger
}

// Body returns the error body for a code, linking to its docs section.
func (e Errors) Body(code, message string) gen.Error {
	return gen.Error{Code: code, Message: message, DocsUrl: e.DocsBaseURL + "#" + code}
}

// Write sends an error response.
func (e Errors) Write(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(e.Body(code, message)); err != nil {
		e.Log.Error("write error body", "err", err)
	}
}

// RequestError handles requests the generated code could not parse.
func (e Errors) RequestError(w http.ResponseWriter, _ *http.Request, err error) {
	e.Write(w, http.StatusBadRequest, CodeBadRequest, err.Error())
}

// ResponseError handles errors returned by handlers. Details are logged, not leaked.
func (e Errors) ResponseError(w http.ResponseWriter, r *http.Request, err error) {
	e.Log.ErrorContext(r.Context(), "handler failed", "err", err, "path", r.URL.Path)
	e.Write(w, http.StatusInternalServerError, CodeInternal, "Internal server error")
}

// NotFound answers unknown routes.
func (e Errors) NotFound(w http.ResponseWriter, _ *http.Request) {
	e.Write(w, http.StatusNotFound, CodeNotFound, "Route not found")
}

// MethodNotAllowed answers known routes called with the wrong method.
func (e Errors) MethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	e.Write(w, http.StatusMethodNotAllowed, CodeMethodNotAllowed, "Method not allowed")
}
