package lexicon

import (
	"context"
	"slices"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// Reader is what the handler needs from the lexicon.
type Reader interface {
	Search(ctx context.Context, language, query string, limit int) ([]Hit, error)
	Get(ctx context.Context, id uuid.UUID) (Entry, bool, error)
	ConceptEntries(ctx context.Context, id uuid.UUID) (ConceptEntries, bool, error)
	FromEnglish(ctx context.Context, query string, languages []string, limit int) ([]LookupHit, error)
	ToEnglish(ctx context.Context, language, query string, limit int) ([]LookupHit, error)
}

// Handler implements the lexicon read operations of the generated strict server.
type Handler struct {
	svc Reader
}

// NewHandler returns a lexicon handler.
func NewHandler(svc Reader) *Handler {
	return &Handler{svc: svc}
}

// SearchEntries implements GET /v1/entries/search.
func (h *Handler) SearchEntries(ctx context.Context, req gen.SearchEntriesRequestObject) (gen.SearchEntriesResponseObject, error) {
	hits, err := h.svc.Search(ctx, req.Params.Language, req.Params.Q, httpx.DerefOr(req.Params.Limit, 20))
	if err != nil {
		return nil, err
	}
	data := make([]gen.EntryHit, len(hits))
	for i, hit := range hits {
		data[i] = gen.EntryHit{Match: gen.Match(hit.Match), Entry: entrySummary(hit.Entry)}
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

// GetEntry implements GET /v1/entries/{id}.
func (h *Handler) GetEntry(ctx context.Context, req gen.GetEntryRequestObject) (gen.GetEntryResponseObject, error) {
	e, ok, err := h.svc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, httpx.ErrNotFound
	}
	senses := make([]gen.Sense, len(e.Senses))
	for i, sn := range e.Senses {
		senses[i] = sense(sn)
	}
	variants := make([]gen.Variant, len(e.Variants))
	for i, v := range e.Variants {
		variants[i] = gen.Variant{Spelling: v.Spelling, Note: httpx.Optional(v.Note)}
	}
	return gen.GetEntry200JSONResponse{
		Id: e.ID, Language: e.Language, Headword: e.Headword, Register: gen.Register(e.Register),
		PartOfSpeech: httpx.Optional(gen.PartOfSpeech(e.PartOfSpeech)), NounClass: httpx.Optional(e.NounClass),
		Plural: httpx.Optional(e.Plural), Root: httpx.Optional(e.Root), Variety: httpx.Optional(e.Variety),
		Variants: variants, Senses: senses, Trust: e.Trust.API(),
	}, nil
}

func sense(s Sense) gen.Sense {
	examples := make([]gen.Example, len(s.Examples))
	for i, x := range s.Examples {
		examples[i] = gen.Example{Id: x.ID, Text: x.Text, TranslationEn: x.TranslationEn, Trust: x.Trust.API()}
	}
	out := gen.Sense{Id: s.ID, Gloss: s.Gloss, Definition: httpx.Optional(s.Definition), Examples: examples, Trust: s.Trust.API()}
	if s.ConceptID.Valid {
		out.ConceptId = &s.ConceptID.UUID
	}
	return out
}

// GetConcept implements GET /v1/concepts/{id}.
func (h *Handler) GetConcept(ctx context.Context, req gen.GetConceptRequestObject) (gen.GetConceptResponseObject, error) {
	c, ok, err := h.svc.ConceptEntries(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, httpx.ErrNotFound
	}
	entries := make([]gen.EntrySummary, len(c.Entries))
	for i, e := range c.Entries {
		entries[i] = entrySummary(e)
	}
	return gen.GetConcept200JSONResponse{
		Id: c.ID, GlossEn: c.GlossEn, Definition: httpx.Optional(c.Definition), Domain: c.Domain, Core: c.Core, Entries: entries,
	}, nil
}

// english is the code lookups translate to and from.
const english = "eng"

// Lookup implements GET /v1/lookup.
func (h *Handler) Lookup(ctx context.Context, req gen.LookupRequestObject) (gen.LookupResponseObject, error) {
	p := req.Params
	to := httpx.Deref(p.To)
	limit := httpx.DerefOr(p.Limit, 5)
	var hits []LookupHit
	var err error
	switch {
	case p.From == english:
		hits, err = h.svc.FromEnglish(ctx, p.Q, to, limit)
	case len(to) > 0 && !slices.Equal(to, []string{english}):
		return nil, httpx.BadRequest("From a Zambian language, lookups translate only to eng.")
	default:
		hits, err = h.svc.ToEnglish(ctx, p.From, p.Q, limit)
	}
	if err != nil {
		return nil, err
	}
	data := make([]gen.LookupHit, len(hits))
	for i, hit := range hits {
		entries := make([]gen.EntrySummary, len(hit.Entries))
		for j, e := range hit.Entries {
			entries[j] = entrySummary(e)
		}
		data[i] = gen.LookupHit{Match: gen.Match(hit.Match), Entries: entries}
		data[i].Concept.Id, data[i].Concept.GlossEn = hit.ConceptID, hit.GlossEn
	}
	return gen.Lookup200JSONResponse{Data: data}, nil
}
