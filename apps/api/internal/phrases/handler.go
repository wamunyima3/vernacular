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
	Match(ctx context.Context, language, text string, fromEnglish bool) ([]Phrase, error)
}

// english is the other side of every translation for now.
const english = "eng"

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
		data[i] = phrase(ph)
	}
	return gen.GetPhrasePack200JSONResponse{Pack: p.Pack, Language: p.Lang, Data: data}, nil
}

// Translate implements POST /v1/translate: exact match against verified
// phrases, between English and one Zambian language (TRN-2).
func (h *Handler) Translate(ctx context.Context, req gen.TranslateRequestObject) (gen.TranslateResponseObject, error) {
	b := req.Body
	if (b.From == english) == (b.To == english) {
		return nil, httpx.BadRequest("Translate between English (eng) and one Zambian language.")
	}
	fromEnglish, language := b.From == english, b.From
	if fromEnglish {
		language = b.To
	}
	phrases, err := h.svc.Match(ctx, language, b.Text, fromEnglish)
	if err != nil {
		return nil, err
	}
	data := make([]gen.Translation, len(phrases))
	for i, ph := range phrases {
		translation := ph.SourceEn
		if fromEnglish {
			translation = ph.Text
		}
		data[i] = gen.Translation{Method: gen.TranslationMethodVerifiedPhrase, Translation: translation, Phrase: phrase(ph)}
	}
	return gen.Translate200JSONResponse{Data: data}, nil
}

func phrase(ph Phrase) gen.Phrase {
	return gen.Phrase{
		Id: ph.ID, Text: ph.Text, SourceEn: ph.SourceEn, Notes: httpx.Optional(ph.Notes),
		ConceptId: httpx.OptionalUUID(ph.ConceptID), Trust: ph.Trust.API(),
	}
}
