package review

import (
	"context"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

// Handler implements the review operations of the generated strict server.
type Handler struct {
	svc *Service
}

// NewHandler returns a review handler.
func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// GetNextReview implements GET /v1/review/next.
func (h *Handler) GetNextReview(ctx context.Context, req gen.GetNextReviewRequestObject) (gen.GetNextReviewResponseObject, error) {
	card, ok, err := h.svc.Next(ctx, req.Params.Language, httpx.Deref(req.Params.Skip))
	if err != nil {
		return nil, err
	}
	if !ok {
		return gen.GetNextReview204Response{}, nil
	}
	candidates := make([]gen.ReviewCandidate, len(card.Candidates))
	for i, c := range card.Candidates {
		candidates[i] = candidate(c)
	}
	return gen.GetNextReview200JSONResponse{
		Language: card.Language, ConceptId: card.ConceptID, GlossEn: card.GlossEn,
		Priority: gen.ReviewPriority(card.Priority), Candidates: candidates,
	}, nil
}

// CreateReviews implements POST /v1/reviews.
func (h *Handler) CreateReviews(ctx context.Context, req gen.CreateReviewsRequestObject) (gen.CreateReviewsResponseObject, error) {
	decisions := make([]Decision, len(req.Body.Decisions))
	for i, d := range req.Body.Decisions {
		decisions[i] = Decision{SenseID: d.SenseId, Decision: dbgen.ReviewDecision(d.Decision)}
	}
	var correction *Correction
	if c := req.Body.Correction; c != nil {
		correction = &Correction{
			Headword: c.Headword, PartOfSpeech: string(httpx.Deref(c.PartOfSpeech)), NounClass: httpx.Deref(c.NounClass),
		}
	}
	recorded, err := h.svc.Decide(ctx, decisions, httpx.Deref(req.Body.Comment), correction)
	if err != nil {
		return nil, err
	}
	out := make([]gen.Review, len(recorded))
	for i, r := range recorded {
		out[i] = gen.Review{
			Id: r.ID, SenseId: r.SenseID, Decision: gen.ReviewDecision(r.Decision.Decision), CorrectionSenseId: r.CorrectionSenseID,
		}
	}
	return gen.CreateReviews201JSONResponse{Data: out}, nil
}

// DeleteReviews implements DELETE /v1/reviews.
func (h *Handler) DeleteReviews(ctx context.Context, req gen.DeleteReviewsRequestObject) (gen.DeleteReviewsResponseObject, error) {
	if err := h.svc.Undo(ctx, req.Params.Ids); err != nil {
		return nil, err
	}
	return gen.DeleteReviews204Response{}, nil
}

func candidate(c lexicon.SenseCard) gen.ReviewCandidate {
	return gen.ReviewCandidate{
		SenseId: c.ID, Headword: c.Headword, PartOfSpeech: httpx.Optional(gen.PartOfSpeech(c.PartOfSpeech)),
		NounClass: httpx.Optional(c.NounClass), Plural: httpx.Optional(c.Plural),
	}
}

// ListDisputes implements GET /v1/disputes.
func (h *Handler) ListDisputes(ctx context.Context, _ gen.ListDisputesRequestObject) (gen.ListDisputesResponseObject, error) {
	disputes, err := h.svc.Disputes(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]gen.Dispute, len(disputes))
	for i, d := range disputes {
		reviews := make([]gen.DisputeReview, len(d.Reviews))
		for j, r := range d.Reviews {
			reviews[j] = gen.DisputeReview{
				Decision: gen.ReviewDecision(r.Decision), Comment: r.Comment,
				Correction: httpx.Optional(r.Correction), CreatedAt: r.CreatedAt,
			}
		}
		out[i] = gen.Dispute{
			Language: d.Language, GlossEn: d.Candidate.GlossEn, Candidate: candidate(d.Candidate),
			DisputedSince: d.DisputedSince, Reviews: reviews,
		}
	}
	return gen.ListDisputes200JSONResponse{Disputes: out}, nil
}

// ResolveDispute implements POST /v1/disputes/{sense_id}/resolution.
func (h *Handler) ResolveDispute(ctx context.Context, req gen.ResolveDisputeRequestObject) (gen.ResolveDisputeResponseObject, error) {
	outcome := dbgen.ResolutionOutcome(req.Body.Outcome)
	status, err := h.svc.Resolve(ctx, req.SenseId, outcome, req.Body.Note)
	if err != nil {
		return nil, err
	}
	return gen.ResolveDispute200JSONResponse{SenseId: req.SenseId, Outcome: req.Body.Outcome, Status: gen.TrustStatus(status)}, nil
}

// GetProviderQuality implements GET /v1/provider-quality.
func (h *Handler) GetProviderQuality(ctx context.Context, _ gen.GetProviderQualityRequestObject) (gen.GetProviderQualityResponseObject, error) {
	rows, err := h.svc.ProviderQuality(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]gen.ProviderQuality, len(rows))
	for i, r := range rows {
		out[i] = gen.ProviderQuality{
			Provider: r.Provider, Language: r.LanguageCode, Decided: int(r.Decided), Accepted: int(r.Accepted), Corrected: int(r.Corrected),
		}
	}
	return gen.GetProviderQuality200JSONResponse{Providers: out}, nil
}

// ListCredits implements GET /v1/credits.
func (h *Handler) ListCredits(ctx context.Context, req gen.ListCreditsRequestObject) (gen.ListCreditsResponseObject, error) {
	rows, err := h.svc.Credits(ctx, req.Params.Language)
	if err != nil {
		return nil, err
	}
	out := make([]gen.Credit, len(rows))
	for i, r := range rows {
		out[i] = gen.Credit{Name: r.Name, Words: int(r.Words), Reviews: int(r.Reviews)}
	}
	return gen.ListCredits200JSONResponse{Credits: out}, nil
}

// CreateReport implements POST /v1/reports.
func (h *Handler) CreateReport(ctx context.Context, req gen.CreateReportRequestObject) (gen.CreateReportResponseObject, error) {
	id, err := h.svc.Report(ctx, req.Body.SenseId, req.Body.Reason)
	if err != nil {
		return nil, err
	}
	return gen.CreateReport201JSONResponse{Id: id}, nil
}

// ListReports implements GET /v1/reports.
func (h *Handler) ListReports(ctx context.Context, req gen.ListReportsRequestObject) (gen.ListReportsResponseObject, error) {
	rows, err := h.svc.Reports(ctx, req.Params.Language)
	if err != nil {
		return nil, err
	}
	data := make([]gen.Report, len(rows))
	for i, r := range rows {
		data[i] = gen.Report{
			Id: r.ID, SenseId: r.ItemID, Language: r.LanguageCode, Headword: r.Headword, Gloss: r.Gloss,
			Status: gen.TrustStatus(r.Status), Reason: r.Reason, CreatedAt: r.CreatedAt,
		}
	}
	return gen.ListReports200JSONResponse{Data: data}, nil
}

// DecideReport implements POST /v1/reports/{id}/decision.
func (h *Handler) DecideReport(ctx context.Context, req gen.DecideReportRequestObject) (gen.DecideReportResponseObject, error) {
	status, err := h.svc.DecideReport(ctx, req.Id, req.Body.Confirm, req.Body.Note)
	if err != nil {
		return nil, err
	}
	return gen.DecideReport200JSONResponse{Status: gen.TrustStatus(status)}, nil
}
