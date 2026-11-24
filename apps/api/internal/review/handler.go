package review

import (
	"context"
	"time"

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
		Language: card.Language, ConceptId: card.ConceptID, GlossEn: card.GlossEn, Kind: gen.ReviewCardKind(card.Kind),
		Priority: gen.ReviewPriority(card.Priority), Candidates: candidates,
	}, nil
}

// CreateReviews implements POST /v1/reviews.
func (h *Handler) CreateReviews(ctx context.Context, req gen.CreateReviewsRequestObject) (gen.CreateReviewsResponseObject, error) {
	decisions := make([]Decision, len(req.Body.Decisions))
	for i, d := range req.Body.Decisions {
		decisions[i] = Decision{ItemID: d.ItemId, Decision: dbgen.ReviewDecision(d.Decision)}
	}
	var correction *Correction
	if c := req.Body.Correction; c != nil {
		correction = &Correction{
			Headword: c.Headword, PartOfSpeech: string(httpx.Deref(c.PartOfSpeech)), NounClass: httpx.Deref(c.NounClass),
			TranslationEn: httpx.Deref(c.TranslationEn),
		}
	}
	recorded, err := h.svc.Decide(ctx, decisions, httpx.Deref(req.Body.Comment), correction)
	if err != nil {
		return nil, err
	}
	out := make([]gen.Review, len(recorded))
	for i, r := range recorded {
		out[i] = gen.Review{
			Id: r.ID, ItemId: r.ItemID, Decision: gen.ReviewDecision(r.Decision.Decision), CorrectionItemId: r.CorrectionItemID,
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

func candidate(c Candidate) gen.ReviewCandidate {
	return gen.ReviewCandidate{
		ItemId: c.ID, Headword: c.Text, PartOfSpeech: httpx.Optional(gen.PartOfSpeech(c.PartOfSpeech)),
		NounClass: httpx.Optional(c.NounClass), Plural: httpx.Optional(c.Plural), Notes: httpx.Optional(c.Notes),
		Word: httpx.Optional(c.Word), TranslationEn: httpx.Optional(c.TranslationEn), Corrects: correctedMeaning(c.Corrects),
	}
}

func correctedMeaning(c *lexicon.Corrected) *gen.CorrectedMeaning {
	if c == nil {
		return nil
	}
	return &gen.CorrectedMeaning{ItemId: c.ID, Headword: c.Headword, Gloss: c.Gloss}
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
			Language: d.Language, GlossEn: d.Candidate.English, Candidate: candidate(d.Candidate),
			DisputedSince: d.DisputedSince, Reviews: reviews,
		}
	}
	return gen.ListDisputes200JSONResponse{Disputes: out}, nil
}

// ResolveDispute implements POST /v1/disputes/{item_id}/resolution.
func (h *Handler) ResolveDispute(ctx context.Context, req gen.ResolveDisputeRequestObject) (gen.ResolveDisputeResponseObject, error) {
	outcome := dbgen.ResolutionOutcome(req.Body.Outcome)
	status, err := h.svc.Resolve(ctx, req.ItemId, outcome, req.Body.Note)
	if err != nil {
		return nil, err
	}
	return gen.ResolveDispute200JSONResponse{ItemId: req.ItemId, Outcome: req.Body.Outcome, Status: gen.TrustStatus(status)}, nil
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

// SampleRechecks implements POST /v1/rechecks.
func (h *Handler) SampleRechecks(ctx context.Context, req gen.SampleRechecksRequestObject) (gen.SampleRechecksResponseObject, error) {
	quarter := httpx.DerefOr(req.Body.Quarter, Quarter(time.Now()))
	n, err := h.svc.Sample(ctx, quarter)
	if err != nil {
		return nil, err
	}
	return gen.SampleRechecks200JSONResponse{Quarter: quarter, Sampled: int(n)}, nil
}

// GetNextRecheck implements GET /v1/rechecks/next.
func (h *Handler) GetNextRecheck(ctx context.Context, req gen.GetNextRecheckRequestObject) (gen.GetNextRecheckResponseObject, error) {
	r, ok, err := h.svc.NextRecheck(ctx, req.Params.Language, httpx.Deref(req.Params.Skip))
	if err != nil {
		return nil, err
	}
	if !ok {
		return gen.GetNextRecheck204Response{}, nil
	}
	return gen.GetNextRecheck200JSONResponse{Quarter: r.Quarter, Language: req.Params.Language, Candidate: candidate(r.Candidate)}, nil
}

// RecordRecheck implements POST /v1/rechecks/{item_id}.
func (h *Handler) RecordRecheck(ctx context.Context, req gen.RecordRecheckRequestObject) (gen.RecordRecheckResponseObject, error) {
	c, err := h.svc.RecordRecheck(ctx, req.ItemId, dbgen.RecheckOutcome(req.Body.Outcome), httpx.Deref(req.Body.Note))
	if err != nil {
		return nil, err
	}
	return gen.RecordRecheck200JSONResponse{ItemId: req.ItemId, Quarter: c.Quarter, Outcome: req.Body.Outcome, Status: gen.TrustStatus(c.Status)}, nil
}

// ListRecheckBatches implements GET /v1/rechecks/batches.
func (h *Handler) ListRecheckBatches(ctx context.Context, req gen.ListRecheckBatchesRequestObject) (gen.ListRecheckBatchesResponseObject, error) {
	quarter := httpx.DerefOr(req.Params.Quarter, Quarter(time.Now()))
	rows, err := h.svc.Batches(ctx, quarter)
	if err != nil {
		return nil, err
	}
	out := make([]gen.RecheckBatch, len(rows))
	for i, r := range rows {
		out[i] = gen.RecheckBatch{
			Language: r.LanguageCode, SourceId: r.SourceID, Source: r.Source, Model: r.Model,
			Sampled: int(r.Sampled), Checked: int(r.Checked), Failed: int(r.Failed), Fails: Fails(int(r.Checked), int(r.Failed)),
			SentBackAt: httpx.OptionalTime(r.SentBackAt), SentBackItems: httpx.OptionalInt(r.SentBackItems),
		}
	}
	return gen.ListRecheckBatches200JSONResponse{Quarter: quarter, Batches: out}, nil
}

// SendBackRecheckBatch implements POST /v1/rechecks/batches/send-back.
func (h *Handler) SendBackRecheckBatch(ctx context.Context, req gen.SendBackRecheckBatchRequestObject) (gen.SendBackRecheckBatchResponseObject, error) {
	in := req.Body
	n, err := h.svc.SendBackBatch(ctx, in.Quarter, Batch{Language: in.Language, SourceID: in.SourceId, Model: in.Model})
	if err != nil {
		return nil, err
	}
	return gen.SendBackRecheckBatch200JSONResponse{Items: n}, nil
}
