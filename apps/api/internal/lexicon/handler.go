package lexicon

import (
	"context"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// Searcher is what the handler needs from the lexicon.
type Searcher interface {
	Search(ctx context.Context, language, query string, limit int) ([]Hit, error)
}

// Handler implements the lexicon read operations of the generated strict server.
type Handler struct {
	svc Searcher
}

// NewHandler returns a lexicon handler.
func NewHandler(svc Searcher) *Handler {
	return &Handler{svc: svc}
}

// SearchEntries implements GET /v1/entries/search.
func (h *Handler) SearchEntries(ctx context.Context, req gen.SearchEntriesRequestObject) (gen.SearchEntriesResponseObject, error) {
	limit := 20
	if req.Params.Limit != nil {
		limit = *req.Params.Limit
	}
	hits, err := h.svc.Search(ctx, req.Params.Language, req.Params.Q, limit)
	if err != nil {
		return nil, err
	}
	data := make([]gen.EntryHit, len(hits))
	for i, hit := range hits {
		data[i] = gen.EntryHit{Match: gen.EntryHitMatch(hit.Match), Entry: entrySummary(hit.Entry)}
	}
	return gen.SearchEntries200JSONResponse{Data: data}, nil
}

func entrySummary(e Entry) gen.EntrySummary {
	senses := make([]gen.SenseSummary, len(e.Senses))
	for i, s := range e.Senses {
		senses[i] = gen.SenseSummary{Id: s.ID, Gloss: s.Gloss, Trust: s.Trust.API()}
	}
	return gen.EntrySummary{
		Id: e.ID, Language: e.Language, Headword: e.Headword, Register: gen.Register(e.Register),
		PartOfSpeech: httpx.Optional(gen.PartOfSpeech(e.PartOfSpeech)), NounClass: httpx.Optional(e.NounClass),
		Plural: httpx.Optional(e.Plural), Root: httpx.Optional(e.Root), Variety: httpx.Optional(e.Variety),
		Senses: senses, Trust: e.Trust.API(),
	}
}
