// Package httpx holds HTTP concerns shared by every module: the uniform error
// body and the router middleware stack.
package httpx

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
)

// Errors builds the uniform error body {code, message, docs_url} used by every
// response. Codes come from the contract (gen.ErrorCode), so Go, TS and the
// docs page share one list.
type Errors struct {
	DocsBaseURL string
	Log         *slog.Logger
}

// Body returns the error body for a code, linking to its docs section.
func (e Errors) Body(code gen.ErrorCode, message string) gen.Error {
	return gen.Error{Code: code, Message: message, DocsUrl: e.DocsBaseURL + "#" + string(code)}
}

// Write sends an error response.
func (e Errors) Write(w http.ResponseWriter, status int, code gen.ErrorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(e.Body(code, message)); err != nil {
		e.Log.Error("write error body", "err", err)
	}
}

// RequestError handles requests the generated code could not parse.
func (e Errors) RequestError(w http.ResponseWriter, _ *http.Request, err error) {
	e.Write(w, http.StatusBadRequest, gen.ErrorCodeBadRequest, err.Error())
}

// Problem is an error a handler returns to answer with a specific status and
// code. Any other error becomes a logged 500.
type Problem struct {
	Status  int
	Code    gen.ErrorCode
	Message string
}

func (p *Problem) Error() string { return p.Message }

// Common problems. Compare with errors.Is.
var (
	ErrUnauthorized = &Problem{http.StatusUnauthorized, gen.ErrorCodeUnauthorized, "Sign in to use this endpoint"}
	ErrForbidden    = &Problem{http.StatusForbidden, gen.ErrorCodeForbidden, "Your account is not allowed to do this"}
	ErrNotFound     = &Problem{http.StatusNotFound, gen.ErrorCodeNotFound, "Not found"}
)

// BadRequest is a 400 with a message for the caller.
func BadRequest(msg string) error {
	return &Problem{http.StatusBadRequest, gen.ErrorCodeBadRequest, msg}
}

// Conflict is a 409 with a message for the caller.
func Conflict(msg string) error {
	return &Problem{http.StatusConflict, gen.ErrorCodeConflict, msg}
}

// TooLarge is a 413 with a message for the caller.
func TooLarge(msg string) error {
	return &Problem{http.StatusRequestEntityTooLarge, gen.ErrorCodePayloadTooLarge, msg}
}

// ResponseError handles errors returned by handlers: a Problem is sent as is,
// and anything else is logged and returned as a bare 500.
func (e Errors) ResponseError(w http.ResponseWriter, r *http.Request, err error) {
	if p, ok := errors.AsType[*Problem](err); ok {
		e.Write(w, p.Status, p.Code, p.Message)
		return
	}
	e.Log.ErrorContext(r.Context(), "handler failed", "err", err, "path", r.URL.Path)
	e.Write(w, http.StatusInternalServerError, gen.ErrorCodeInternal, "Internal server error")
}

// NotFound answers unknown routes.
func (e Errors) NotFound(w http.ResponseWriter, _ *http.Request) {
	e.Write(w, http.StatusNotFound, gen.ErrorCodeNotFound, "Route not found")
}

// MethodNotAllowed answers known routes called with the wrong method.
func (e Errors) MethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	e.Write(w, http.StatusMethodNotAllowed, gen.ErrorCodeMethodNotAllowed, "Method not allowed")
}
