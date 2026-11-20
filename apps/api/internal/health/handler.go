// Package health answers the liveness and readiness probe.
package health

import (
	"context"
	"time"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// Pinger reports whether a dependency answers. *pgxpool.Pool satisfies it.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Handler implements the health operations of the generated strict server.
type Handler struct {
	db     Pinger
	errors httpx.Errors
}

// NewHandler returns a health handler that checks the database.
func NewHandler(db Pinger, errors httpx.Errors) *Handler {
	return &Handler{db: db, errors: errors}
}

// GetHealth implements GET /healthz.
func (h *Handler) GetHealth(ctx context.Context, _ gen.GetHealthRequestObject) (gen.GetHealthResponseObject, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := h.db.Ping(ctx); err != nil {
		body := h.errors.Body(gen.ErrorCodeUnavailable, "Database unreachable")
		return gen.GetHealth503JSONResponse{ErrorJSONResponse: gen.ErrorJSONResponse(body)}, nil
	}
	return gen.GetHealth200JSONResponse{Status: "ok"}, nil
}
