package review

import (
	"context"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
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
		candidates[i] = gen.ReviewCandidate{
			SenseId: c.ID, Headword: c.Headword, PartOfSpeech: httpx.Optional(gen.PartOfSpeech(c.PartOfSpeech)),
			NounClass: httpx.Optional(c.NounClass), Plural: httpx.Optional(c.Plural),
		}
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
	recorded, err := h.svc.Decide(ctx, decisions, httpx.Deref(req.Body.Comment))
	if err != nil {
		return nil, err
	}
	out := make([]gen.Review, len(recorded))
	for i, r := range recorded {
		out[i] = gen.Review{Id: r.ID, SenseId: r.SenseID, Decision: gen.ReviewDecision(r.Decision.Decision)}
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
