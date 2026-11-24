package contributions

import (
	"context"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
)

// Handler implements the contributions operations of the generated strict server.
type Handler struct {
	svc *Service
}

// NewHandler returns a contributions handler.
func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// CreateContribution implements POST /v1/contributions. The contract has
// already checked the shape, the part of speech and agreement to the licence.
func (h *Handler) CreateContribution(ctx context.Context, req gen.CreateContributionRequestObject) (gen.CreateContributionResponseObject, error) {
	b := req.Body
	in := Input{
		Language: b.Language, Headword: b.Headword, PartOfSpeech: string(b.PartOfSpeech), GlossEn: b.GlossEn,
	}
	if b.NounClass != nil {
		in.NounClass = *b.NounClass
	}
	if b.Corrects != nil {
		in.Corrects = uuid.NullUUID{UUID: *b.Corrects, Valid: true}
	}
	if b.Example != nil {
		in.Example = &Example{Text: b.Example.Text, TranslationEn: b.Example.TranslationEn}
	}
	r, err := h.svc.Propose(ctx, in)
	if err != nil {
		return nil, err
	}
	return gen.CreateContribution201JSONResponse{
		EntryId: r.EntryID, SenseId: r.SenseID, ExampleId: r.ExampleID, NewSense: r.NewSense, Trust: r.Trust.API(),
	}, nil
}
