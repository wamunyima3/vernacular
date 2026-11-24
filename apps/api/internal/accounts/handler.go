package accounts

import (
	"context"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
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
		out[i] = gen.Reviewer{
			UserId: r.ID, Name: r.Name, Level: gen.ReviewerLevel(r.Level), Decided: int(r.Decided), Agreed: int(r.Agreed),
			Suspended: r.Suspended, RetrainedAt: httpx.OptionalTime(r.RetrainedAt),
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

func key(k dbgen.ListKeysRow) gen.Key {
	return gen.Key{
		Id: k.ID, Name: k.Name, Prefix: k.Prefix, CreatedAt: k.CreatedAt,
		LastUsedAt: httpx.OptionalTime(k.LastUsedAt), RevokedAt: httpx.OptionalTime(k.RevokedAt),
	}
}

// ListKeys implements GET /v1/keys.
func (h *Handler) ListKeys(ctx context.Context, _ gen.ListKeysRequestObject) (gen.ListKeysResponseObject, error) {
	rows, err := h.svc.Keys(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]gen.Key, len(rows))
	for i, r := range rows {
		out[i] = key(r)
	}
	return gen.ListKeys200JSONResponse{Data: out}, nil
}

// CreateKey implements POST /v1/keys.
func (h *Handler) CreateKey(ctx context.Context, req gen.CreateKeyRequestObject) (gen.CreateKeyResponseObject, error) {
	k, err := h.svc.CreateKey(ctx, req.Body.Name)
	if err != nil {
		return nil, err
	}
	return gen.CreateKey201JSONResponse{Key: key(k.ListKeysRow), Secret: k.Secret}, nil
}

// RevokeKey implements DELETE /v1/keys/{id}.
func (h *Handler) RevokeKey(ctx context.Context, req gen.RevokeKeyRequestObject) (gen.RevokeKeyResponseObject, error) {
	if err := h.svc.RevokeKey(ctx, req.Id); err != nil {
		return nil, err
	}
	return gen.RevokeKey204Response{}, nil
}

// StartQualification implements POST /v1/qualifications.
func (h *Handler) StartQualification(ctx context.Context, req gen.StartQualificationRequestObject) (gen.StartQualificationResponseObject, error) {
	t, err := h.svc.StartTest(ctx, req.Body.Language)
	if err != nil {
		return nil, err
	}
	questions := make([]gen.QualificationQuestion, len(t.Questions))
	for i, q := range t.Questions {
		options := make([]gen.QualificationOption, len(q.Options))
		for j, o := range q.Options {
			options[j] = gen.QualificationOption{Id: o.ID, Headword: o.Headword}
		}
		questions[i] = gen.QualificationQuestion{GlossEn: q.GlossEn, Options: options}
	}
	return gen.StartQualification201JSONResponse{Id: t.ID, Language: t.Language, Questions: questions}, nil
}

// AnswerQualification implements POST /v1/qualifications/{id}/answers.
func (h *Handler) AnswerQualification(ctx context.Context, req gen.AnswerQualificationRequestObject) (gen.AnswerQualificationResponseObject, error) {
	r, err := h.svc.AnswerTest(ctx, req.Id, req.Body.Answers)
	if err != nil {
		return nil, err
	}
	return gen.AnswerQualification200JSONResponse{Id: req.Id, Score: r.Score, Total: r.Total, Passed: r.Passed}, nil
}

// ConfirmQualification implements POST /v1/qualifications/{id}/confirmation.
func (h *Handler) ConfirmQualification(ctx context.Context, req gen.ConfirmQualificationRequestObject) (gen.ConfirmQualificationResponseObject, error) {
	if err := h.svc.ConfirmTest(ctx, req.Id); err != nil {
		return nil, err
	}
	return gen.ConfirmQualification204Response{}, nil
}

// ListQualifications implements GET /v1/qualifications.
func (h *Handler) ListQualifications(ctx context.Context, _ gen.ListQualificationsRequestObject) (gen.ListQualificationsResponseObject, error) {
	rows, err := h.svc.PendingTests(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]gen.PendingQualification, len(rows))
	for i, r := range rows {
		out[i] = gen.PendingQualification{
			Id: r.ID, UserId: r.UserID, Name: r.Name, Language: r.LanguageCode,
			Score: int(r.Score.Int32), Total: int(r.Total), AnsweredAt: r.AnsweredAt.Time,
		}
	}
	return gen.ListQualifications200JSONResponse{Qualifications: out}, nil
}
