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

// Store reads phrase packs and records phrase candidates.
type Store struct {
	q      *dbgen.Queries
	items  *items.Store
	locked map[string]bool // languages this transaction holds the candidate lock for
}

// NewStore wraps a pool or transaction.
func NewStore(db dbgen.DBTX) *Store {
	return &Store{q: dbgen.New(db), items: items.NewStore(db), locked: map[string]bool{}}
}

// Pack reads a pack's phrases in a language at or above minStatus.
func (s *Store) Pack(ctx context.Context, language, pack, minStatus string) ([]Phrase, error) {
	rows, err := s.q.ListPack(ctx, dbgen.ListPackParams{Language: language, Pack: pack, MinStatus: dbgen.TrustStatus(minStatus)})
	if err != nil {
		return nil, fmt.Errorf("list pack %s/%s: %w", language, pack, err)
	}
	return s.withTrust(ctx, rows)
}

// Match finds verified phrases in a language whose text, or English source
// when fromEnglish, equals the text once folded (TRN-2).
func (s *Store) Match(ctx context.Context, language, text string, fromEnglish bool) ([]Phrase, error) {
	// Two queries rather than one with a CASE, so each uses its exact-match index.
	var rows []dbgen.ListPackRow
	if fromEnglish {
		en, err := s.q.MatchEnglish(ctx, dbgen.MatchEnglishParams{Language: language, Text: text})
		if err != nil {
			return nil, fmt.Errorf("match phrases: %w", err)
		}
		rows = convert(en)
	} else {
		own, err := s.q.MatchText(ctx, dbgen.MatchTextParams{Language: language, Text: text})
		if err != nil {
			return nil, fmt.Errorf("match phrases: %w", err)
		}
		rows = convert(own)
	}
	return s.withTrust(ctx, rows)
}

// withTrust adds the trust fields to phrase rows.
func (s *Store) withTrust(ctx context.Context, rows []dbgen.ListPackRow) ([]Phrase, error) {
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

// convert turns rows of queries selecting the pack's columns into pack rows.
func convert[T ~struct {
	ID        uuid.UUID
	Text      string
	SourceEn  string
	Notes     string
	ConceptID uuid.NullUUID
}](rows []T) []dbgen.ListPackRow {
	out := make([]dbgen.ListPackRow, len(rows))
	for i, r := range rows {
		out[i] = dbgen.ListPackRow(r)
	}
	return out
}
