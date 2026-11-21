// Package phrases owns curated phrase packs (PHR-1, architecture §5).
package phrases

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/phrases/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
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
	if len(rows) == 0 {
		return nil, nil
	}
	ids := database.IDs(rows, func(r dbgen.ListPackRow) uuid.UUID { return r.ID })
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
