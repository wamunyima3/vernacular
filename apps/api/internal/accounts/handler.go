package accounts

import (
	"context"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
)

// Handler implements the account operations of the generated strict server.
type Handler struct {
	svc *Service
}

// NewHandler returns an accounts handler.
func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// GetMe implements GET /v1/me.
func (h *Handler) GetMe(ctx context.Context, _ gen.GetMeRequestObject) (gen.GetMeResponseObject, error) {
	p, err := Require(ctx, nil)
	if err != nil {
		return nil, err
	}
	approvals := make([]gen.Approval, len(p.Approvals))
	for i, a := range p.Approvals {
		approvals[i] = gen.Approval{Language: a.Language, Level: a.Level, Suspended: a.Suspended}
	}
	return gen.GetMe200JSONResponse{
		Id: p.ID, Name: p.Name, Email: openapi_types.Email(p.Email), Role: p.Role, Approvals: approvals, Credited: p.Credited,
	}, nil
}

// ListReviewers implements GET /v1/reviewers.
func (h *Handler) ListReviewers(ctx context.Context, req gen.ListReviewersRequestObject) (gen.ListReviewersResponseObject, error) {
	rows, err := h.svc.Reviewers(ctx, req.Params.Language)
	if err != nil {
		return nil, err
	}
	out := make([]gen.Reviewer, len(rows))
	for i, r := range rows {
		var retrained *time.Time
		if r.RetrainedAt.Valid {
			retrained = &r.RetrainedAt.Time
		}
		out[i] = gen.Reviewer{
			UserId: r.ID, Name: r.Name, Level: gen.ReviewerLevel(r.Level), Decided: int(r.Decided), Agreed: int(r.Agreed),
			Suspended: r.Suspended, RetrainedAt: retrained,
		}
	}
	return gen.ListReviewers200JSONResponse{Reviewers: out}, nil
}

// MarkRetrained implements POST /v1/reviewers/{user_id}/retraining.
func (h *Handler) MarkRetrained(ctx context.Context, req gen.MarkRetrainedRequestObject) (gen.MarkRetrainedResponseObject, error) {
	if err := h.svc.Retrain(ctx, req.UserId, req.Body.Language); err != nil {
		return nil, err
	}
	return gen.MarkRetrained204Response{}, nil
}

// ListApprovals implements GET /v1/approvals.
func (h *Handler) ListApprovals(ctx context.Context, _ gen.ListApprovalsRequestObject) (gen.ListApprovalsResponseObject, error) {
	rows, err := h.svc.Approvals(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]gen.UserApproval, len(rows))
	for i, r := range rows {
		out[i] = gen.UserApproval{
			UserId: r.ID, Name: r.Name, Email: openapi_types.Email(r.Email), Language: r.LanguageCode,
			Level: gen.ReviewerLevel(r.Level), Suspended: r.Suspended,
		}
	}
	return gen.ListApprovals200JSONResponse{Approvals: out}, nil
}

// ApproveReviewer implements PUT /v1/approvals.
func (h *Handler) ApproveReviewer(ctx context.Context, req gen.ApproveReviewerRequestObject) (gen.ApproveReviewerResponseObject, error) {
	in := req.Body
	if err := h.svc.Approve(ctx, string(in.Email), in.Language, in.Level); err != nil {
		return nil, err
	}
	return gen.ApproveReviewer204Response{}, nil
}

// RevokeApproval implements DELETE /v1/approvals/{user_id}/{language}.
func (h *Handler) RevokeApproval(ctx context.Context, req gen.RevokeApprovalRequestObject) (gen.RevokeApprovalResponseObject, error) {
	if err := h.svc.Revoke(ctx, req.UserId, req.Language); err != nil {
		return nil, err
	}
	return gen.RevokeApproval204Response{}, nil
}

// SetCredit implements PUT /v1/me/credit.
func (h *Handler) SetCredit(ctx context.Context, req gen.SetCreditRequestObject) (gen.SetCreditResponseObject, error) {
	if err := h.svc.SetCredit(ctx, req.Body.Credited); err != nil {
		return nil, err
	}
	return gen.SetCredit204Response{}, nil
}
