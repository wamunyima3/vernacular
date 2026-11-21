package phrases

import (
	"context"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// Reader is what the handler needs.
type Reader interface {
	Pack(ctx context.Context, language, pack, minStatus string) ([]Phrase, error)
}

// Handler implements the phrase operations of the generated strict server.
type Handler struct {
	svc Reader
}

// NewHandler returns a phrases handler.
func NewHandler(svc Reader) *Handler {
	return &Handler{svc: svc}
}

// GetPhrasePack implements GET /v1/phrases.
func (h *Handler) GetPhrasePack(ctx context.Context, req gen.GetPhrasePackRequestObject) (gen.GetPhrasePackResponseObject, error) {
	p := req.Params
	phrases, err := h.svc.Pack(ctx, p.Lang, p.Pack, items.Minimum(p.MinStatus))
	if err != nil {
		return nil, err
	}
	data := make([]gen.Phrase, len(phrases))
	for i, ph := range phrases {
		data[i] = gen.Phrase{
			Id: ph.ID, Text: ph.Text, SourceEn: ph.SourceEn, Notes: httpx.Optional(ph.Notes),
			ConceptId: httpx.OptionalUUID(ph.ConceptID), Trust: ph.Trust.API(),
		}
	}
	return gen.GetPhrasePack200JSONResponse{Pack: p.Pack, Language: p.Lang, Data: data}, nil
}
