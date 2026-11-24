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
// agree raise one candidate (AI-3). Give it a transaction. It reports
// whether a new phrase was made.
func (s *Store) Propose(ctx context.Context, language string, line Line, text, notes string, origin items.Origin) (bool, error) {
	text = textnorm.Squash(textnorm.NFC(text))
	if text == "" {
		return false, nil
	}
	if err := s.lock(ctx, language); err != nil {
		return false, err
	}
	id, err := s.q.FindPhrase(ctx, dbgen.FindPhraseParams{Language: language, Pack: line.Pack, SourceEn: line.SourceEn, Text: text})
	if err == nil {
		return false, s.items.Link(ctx, id, origin)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return false, fmt.Errorf("find phrase: %w", err)
	}
	_, err = s.items.Create(ctx, items.KindPhrase, language, origin, func(id uuid.UUID) error {
		return s.q.InsertPhrase(ctx, dbgen.InsertPhraseParams{
			ID: id, Language: language, ConceptID: uuid.NullUUID{UUID: line.ConceptID, Valid: true},
			Pack: line.Pack, Text: text, SourceEn: line.SourceEn, Notes: textnorm.Squash(textnorm.NFC(notes)),
		})
	})
	return err == nil, err
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
