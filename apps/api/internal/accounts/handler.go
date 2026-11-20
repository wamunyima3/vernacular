package accounts

import (
	"context"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
)

// Handler implements the account operations of the generated strict server.
type Handler struct{}

// NewHandler returns an accounts handler.
func NewHandler() *Handler { return &Handler{} }

// GetMe implements GET /v1/me.
func (h *Handler) GetMe(ctx context.Context, _ gen.GetMeRequestObject) (gen.GetMeResponseObject, error) {
	p, err := Require(ctx, nil)
	if err != nil {
		return nil, err
	}
	approvals := make([]gen.Approval, len(p.Approvals))
	for i, a := range p.Approvals {
		approvals[i] = gen.Approval{Language: a.Language, Level: a.Level}
	}
	return gen.GetMe200JSONResponse{
		Id: p.ID, Name: p.Name, Email: openapi_types.Email(p.Email), Role: p.Role, Approvals: approvals,
	}, nil
}
