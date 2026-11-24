package phrases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/phrases/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/textnorm"
)

// Line is one line of the English source list: what a pack should say.
type Line struct {
	ID                    uuid.UUID
	Pack, SourceEn, Notes string
	ConceptID             uuid.UUID
}

// LineIDs lists every line of the source list.
func (s *Store) LineIDs(ctx context.Context) ([]uuid.UUID, error) {
	return s.q.ListLineIDs(ctx)
}

// Line reads one line of the source list.
func (s *Store) Line(ctx context.Context, id uuid.UUID) (Line, error) {
	r, err := s.q.GetListLine(ctx, id)
	if err != nil {
		return Line{}, fmt.Errorf("get phrase line %s: %w", id, err)
	}
	return Line{ID: r.ID, Pack: r.Pack, SourceEn: r.SourceEn, Notes: r.Notes, ConceptID: r.ConceptID}, nil
}

// Propose records text as a phrase for the line in a language, with origin
// as one of its sources. When the same text, folded, was already proposed
// for the line, the origin is linked to that item instead, so providers that
// agree raise one candidate (AI-3). Give it a transaction. It returns the
// phrase, and whether it is new; empty text is no phrase. Text that drops or
// adds a placeholder fails with ErrPlaceholders.
func (s *Store) Propose(ctx context.Context, language string, line Line, text, notes string, origin items.Origin) (uuid.UUID, bool, error) {
	text = textnorm.Squash(textnorm.NFC(text))
	if text == "" {
		return uuid.Nil, false, nil
	}
	if err := checkPlaceholders(line.SourceEn, text); err != nil {
		return uuid.Nil, false, err
	}
	if err := s.lock(ctx, language); err != nil {
		return uuid.Nil, false, err
	}
	id, err := s.q.FindPhrase(ctx, dbgen.FindPhraseParams{Language: language, Pack: line.Pack, SourceEn: line.SourceEn, Text: text})
	if err == nil {
		return id, false, s.items.Link(ctx, id, origin)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, false, fmt.Errorf("find phrase: %w", err)
	}
	id, err = s.items.Create(ctx, items.KindPhrase, language, origin, func(id uuid.UUID) error {
		return s.q.InsertPhrase(ctx, dbgen.InsertPhraseParams{
			ID: id, Language: language, ConceptID: uuid.NullUUID{UUID: line.ConceptID, Valid: true},
			Pack: line.Pack, Text: text, SourceEn: line.SourceEn, Notes: textnorm.Squash(textnorm.NFC(notes)),
		})
	})
	return id, err == nil, err
}

// LineOf is the source-list line a phrase translates: its pack, English and
// concept.
func (s *Store) LineOf(ctx context.Context, phraseID uuid.UUID) (Line, error) {
	r, err := s.q.GetPhraseLine(ctx, phraseID)
	if err != nil {
		return Line{}, fmt.Errorf("get line of phrase %s: %w", phraseID, err)
	}
	return Line{ID: r.ID, Pack: r.Pack, SourceEn: r.SourceEn, Notes: r.Notes, ConceptID: r.ConceptID}, nil
}

// lock takes the language's candidate lock once per transaction.
func (s *Store) lock(ctx context.Context, language string) error {
	if s.locked[language] {
		return nil
	}
	if err := s.q.LockPhrases(ctx, language); err != nil {
		return fmt.Errorf("lock phrases %s: %w", language, err)
	}
	s.locked[language] = true
	return nil
}

// Cards reads the given phrases for review, in the order of ids, without
// trust fields: reviewers judge them blind (CON-3).
func (s *Store) Cards(ctx context.Context, ids []uuid.UUID) ([]Phrase, error) {
	rows, err := s.q.ListPhrasesByID(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("list phrases: %w", err)
	}
	byID := make(map[uuid.UUID]dbgen.ListPackRow, len(rows))
	for _, r := range convert(rows) {
		byID[r.ID] = r
	}
	out := make([]Phrase, 0, len(ids))
	for _, id := range ids {
		if r, ok := byID[id]; ok {
			out = append(out, Phrase{ID: r.ID, Text: r.Text, SourceEn: r.SourceEn, Notes: r.Notes, ConceptID: r.ConceptID})
		}
	}
	return out, nil
}
