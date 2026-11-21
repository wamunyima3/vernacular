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
