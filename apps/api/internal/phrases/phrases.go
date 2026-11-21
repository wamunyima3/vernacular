// Package phrases owns curated phrase packs (PHR-1, architecture §5).
package phrases

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/phrases/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// Phrase is one phrase of a pack, with its trust fields (PRV-1).
type Phrase struct {
	ID                    uuid.UUID
	Text, SourceEn, Notes string
	ConceptID             uuid.NullUUID
	Trust                 items.Trust
}

// Store reads phrase packs.
type Store struct {
	q     *dbgen.Queries
	items *items.Store
}

// NewStore wraps a pool or transaction.
func NewStore(db dbgen.DBTX) *Store {
	return &Store{q: dbgen.New(db), items: items.NewStore(db)}
}

// Pack reads a pack's phrases in a language at or above minStatus.
func (s *Store) Pack(ctx context.Context, language, pack, minStatus string) ([]Phrase, error) {
	rows, err := s.q.ListPack(ctx, dbgen.ListPackParams{Language: language, Pack: pack, MinStatus: dbgen.TrustStatus(minStatus)})
	if err != nil {
		return nil, fmt.Errorf("list pack %s/%s: %w", language, pack, err)
	}
	ids := make([]uuid.UUID, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
	}
	trusts, err := s.items.Trusts(ctx, ids)
	if err != nil {
		return nil, err
	}
	out := make([]Phrase, len(rows))
	for i, r := range rows {
		out[i] = Phrase{ID: r.ID, Text: r.Text, SourceEn: r.SourceEn, Notes: r.Notes, ConceptID: r.ConceptID, Trust: trusts[r.ID]}
	}
	return out, nil
}

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
			Id: ph.ID, Text: ph.Text, SourceEn: ph.SourceEn, Notes: httpx.Optional(ph.Notes), Trust: ph.Trust.API(),
		}
		if ph.ConceptID.Valid {
			data[i].ConceptId = &ph.ConceptID.UUID
		}
	}
	return gen.GetPhrasePack200JSONResponse{Pack: p.Pack, Language: p.Lang, Data: data}, nil
}
