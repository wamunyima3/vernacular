// Package review owns the blind review queue, decisions and, from P1-REV-03,
// promotion to `verified` (CON-2, CON-3, DQ-4, AI-4). The queue reads the
// items table and the review_agreement view; item state goes through items.
package review

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

// Card is one English meaning with every candidate still open to the reviewer.
type Card struct {
	Language   string
	ConceptID  uuid.UUID
	GlossEn    string
	Priority   dbgen.ReviewPriority
	Candidates []lexicon.SenseCard
}

// Decision is a reviewer's verdict on one candidate.
type Decision struct {
	SenseID  uuid.UUID
	Decision dbgen.ReviewDecision
}

// Recorded is a stored decision.
type Recorded struct {
	ID uuid.UUID
	Decision
}

// Lexicon is what the queue needs from the lexicon module.
type Lexicon interface {
	Concept(ctx context.Context, id uuid.UUID) (lexicon.Concept, error)
	SenseCards(ctx context.Context, ids []uuid.UUID) ([]lexicon.SenseCard, error)
}

// Items is what decisions need from the items module.
type Items interface {
	Lock(ctx context.Context, id uuid.UUID) (items.Locked, bool, error)
}

// Bound to the request's transaction.
var (
	newLexicon = func(tx pgx.Tx) Lexicon { return lexicon.NewStore(tx) }
	newItems   = func(tx pgx.Tx) Items { return items.NewStore(tx) }
)

// Service runs the review queue.
type Service struct {
	db database.Beginner
}

// NewService returns a review service.
func NewService(db database.Beginner) *Service { return &Service{db: db} }

// Next returns the next card in a language, or false when nothing is left.
func (s *Service) Next(ctx context.Context, language string, skip []uuid.UUID) (Card, bool, error) {
	me, err := accounts.Require(ctx, func(p accounts.Principal) bool { return p.CanReview(language) })
	if err != nil {
		return Card{}, false, err
	}
	var card Card
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		rows, err := NewStore(tx).nextCard(ctx, language, me.ID, skip)
		if err != nil || len(rows) == 0 {
			return err
		}
		ids := make([]uuid.UUID, len(rows))
		for i, r := range rows {
			ids[i] = r.ItemID
		}
		lex := newLexicon(tx)
		concept, err := lex.Concept(ctx, rows[0].ConceptID)
		if err != nil {
			return err
		}
		cards, err := lex.SenseCards(ctx, ids)
		if err != nil {
			return err
		}
		card = Card{Language: language, ConceptID: concept.ID, GlossEn: concept.GlossEn, Priority: rows[0].Priority, Candidates: cards}
		return nil
	})
	return card, card.Candidates != nil, err
}

// Decide records the reviewer's decisions on candidates from one card, all
// or none.
func (s *Service) Decide(ctx context.Context, decisions []Decision, comment string) ([]Recorded, error) {
	me, err := accounts.Require(ctx, nil)
	if err != nil {
		return nil, err
	}
	out := make([]Recorded, 0, len(decisions))
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		store, its := NewStore(tx), newItems(tx)
		for _, d := range decisions {
			if err := checkTarget(ctx, its, me, d.SenseID); err != nil {
				return err
			}
			id, err := store.create(ctx, d.SenseID, me.ID, d.Decision, comment)
			if errors.Is(err, errAlreadyReviewed) {
				return httpx.Conflict("You have already reviewed " + d.SenseID.String())
			}
			if err != nil {
				return err
			}
			out = append(out, Recorded{ID: id, Decision: d})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// checkTarget locks one item and applies the blind-review rules (CON-2, CON-3).
func checkTarget(ctx context.Context, its Items, me accounts.Principal, id uuid.UUID) error {
	t, ok, err := its.Lock(ctx, id)
	switch {
	case err != nil:
		return err
	case !ok || t.Kind != items.KindSense:
		return httpx.BadRequest("No sense " + id.String())
	case !me.CanReview(t.Language):
		return httpx.ErrForbidden
	case t.AuthorID == me.ID:
		return httpx.Forbidden("You cannot review an item you wrote")
	case !t.Open:
		return httpx.Conflict("This item is " + string(t.Status) + " and no longer in the queue")
	}
	return nil
}
