// Package review owns the blind review queue, decisions and, from P1-REV-03,
// promotion to `verified` (CON-2, CON-3, DQ-4, AI-4). The queue reads the
// items table and the review_agreement view; item state goes through items.
package review

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/phrases"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

// Card is one English meaning with every candidate still open to the reviewer.
type Card struct {
	Language  string
	ConceptID uuid.UUID
	GlossEn   string
	Priority  dbgen.ReviewPriority
	Kind      dbgen.ItemKind
	// The candidates: senses on a word card, phrases on a phrase card.
	Candidates []lexicon.SenseCard
	Phrases    []phrases.Phrase
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
	CorrectionSenseID *uuid.UUID // the new candidate, for corrections
}

// Lexicon is what the queue needs from the lexicon module.
type Lexicon interface {
	Concept(ctx context.Context, id uuid.UUID) (lexicon.Concept, error)
	SenseCards(ctx context.Context, ids []uuid.UUID) ([]lexicon.SenseCard, error)
}

// Phrases is what the queue needs from the phrases module.
type Phrases interface {
	Cards(ctx context.Context, ids []uuid.UUID) ([]phrases.Phrase, error)
	LineOf(ctx context.Context, id uuid.UUID) (phrases.Line, error)
	Propose(ctx context.Context, language string, line phrases.Line, text, notes string, origin items.Origin) (uuid.UUID, bool, error)
}

// Items is what decisions need from the items module.
type Items interface {
	Lock(ctx context.Context, id uuid.UUID) (items.Locked, bool, error)
	SetStatus(ctx context.Context, id uuid.UUID, status items.Status) error
}

// Bound to the request's transaction.
var (
	newLexicon = func(tx pgx.Tx) Lexicon { return lexicon.NewStore(tx) }
	newItems   = func(tx pgx.Tx) Items { return items.NewStore(tx) }
	newPhrases = func(tx pgx.Tx) Phrases { return phrases.NewStore(tx) }
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
		card = Card{Language: language, ConceptID: rows[0].ConceptID, Priority: rows[0].Priority, Kind: rows[0].Kind}
		if card.Kind == dbgen.ItemKindPhrase {
			// A phrase card shows the English as the pack writes it.
			if card.Phrases, err = newPhrases(tx).Cards(ctx, ids); err == nil && len(card.Phrases) > 0 {
				card.GlossEn = card.Phrases[0].SourceEn
			}
			return err
		}
		lex := newLexicon(tx)
		concept, err := lex.Concept(ctx, card.ConceptID)
		if err != nil {
			return err
		}
		card.GlossEn = concept.GlossEn
		card.Candidates, err = lex.SenseCards(ctx, ids)
		return err
	})
	return card, card.Kind != "", err
}

// Decide records the reviewer's decisions on candidates from one card, all
// or none.
func (s *Service) Decide(ctx context.Context, decisions []Decision, comment string, correction *Correction) ([]Recorded, error) {
	me, err := accounts.Require(ctx, nil)
	if err != nil {
		return nil, err
	}
	if corrects := slices.ContainsFunc(decisions, isCorrection); corrects != (correction != nil) {
		return nil, httpx.BadRequest("A correction needs `correct` decisions, and `correct` decisions need a correction")
	}
	// Lock items in one order everywhere, so concurrent decisions and undos cannot deadlock.
	sorted := slices.SortedFunc(slices.Values(decisions), bySense)
	out := make([]Recorded, 0, len(decisions))
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		store, its := NewStore(tx), newItems(tx)
		locked := make([]items.Locked, len(sorted))
		for i, d := range sorted {
			if locked[i], err = checkTarget(ctx, its, me, d.SenseID); err != nil {
				return err
			}
			if locked[i].Kind != locked[0].Kind {
				return httpx.BadRequest("Decide on one card's candidates together: words and phrases are never on one card")
			}
		}
		fix, err := correct(ctx, tx, me.ID, locked, sorted, correction)
		if err != nil {
			return err
		}
		for i, d := range sorted {
			r := Recorded{Decision: d}
			var detail []byte
			if isCorrection(d) {
				r.CorrectionSenseID, detail = &fix.senseID, fix.detail
			}
			if r.ID, err = store.create(ctx, d.SenseID, me.ID, d.Decision, comment, detail); errors.Is(err, errAlreadyReviewed) {
				return httpx.Conflict("You have already reviewed " + d.SenseID.String())
			} else if err != nil {
				return err
			}
			if _, err := promote(ctx, store, its, d.SenseID, locked[i].Origin); err != nil {
				return err
			}
			out = append(out, r)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UndoWindow is how long a reviewer can take back a decision.
const UndoWindow = 10 * time.Minute

// Undo takes back the reviews one decision recorded, all or none, putting
// the items back in the reviewer's queue. Reviews by others count as not found.
func (s *Service) Undo(ctx context.Context, ids []uuid.UUID) error {
	me, err := accounts.Require(ctx, nil)
	if err != nil {
		return err
	}
	return database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		store, its := NewStore(tx), newItems(tx)
		own, err := store.own(ctx, ids, me.ID)
		if err != nil {
			return err
		}
		distinct := make(map[uuid.UUID]bool, len(ids))
		for _, id := range ids {
			distinct[id] = true
		}
		if len(own) != len(distinct) {
			return httpx.ErrNotFound
		}
		targets := make([]uuid.UUID, len(own))
		for i, r := range own {
			if time.Since(r.CreatedAt) > UndoWindow {
				return httpx.Conflict(fmt.Sprintf("Decisions can only be undone within %v", UndoWindow))
			}
			targets[i] = r.ItemID
		}
		slices.SortFunc(targets, compareIDs) // the same lock order as Decide
		locked := make([]items.Locked, len(targets))
		for i, id := range targets {
			if locked[i], _, err = its.Lock(ctx, id); err != nil {
				return err
			}
		}
		if err := store.delete(ctx, ids); err != nil {
			return err
		}
		// Recompute through the promotion rule, so a promoting accept can be undone.
		for i, id := range targets {
			if _, err := promote(ctx, store, its, id, locked[i].Origin); err != nil {
				return err
			}
		}
		return nil
	})
}

func compareIDs(a, b uuid.UUID) int { return bytes.Compare(a[:], b[:]) }

func bySense(a, b Decision) int { return compareIDs(a.SenseID, b.SenseID) }

// correct records the reviewer's own word or phrase when the decisions
// include corrections; nil when they do not.
func correct(ctx context.Context, tx pgx.Tx, reviewerID string, locked []items.Locked, sorted []Decision, c *Correction) (*replacement, error) {
	if c == nil {
		return nil, nil
	}
	if locked[0].Kind == items.KindPhrase {
		return replacePhrase(ctx, tx, reviewerID, locked[0].Language, sorted, *c)
	}
	return replace(ctx, tx, reviewerID, locked[0].Language, sorted, *c)
}

// checkTarget locks one item and applies the blind-review rules (CON-2, CON-3).
func checkTarget(ctx context.Context, its Items, me accounts.Principal, id uuid.UUID) (items.Locked, error) {
	t, ok, err := its.Lock(ctx, id)
	switch {
	case err != nil:
		return t, err
	case !ok || (t.Kind != items.KindSense && t.Kind != items.KindPhrase):
		return t, httpx.BadRequest("No candidate " + id.String())
	case !me.CanReview(t.Language):
		return t, httpx.ErrForbidden
	case t.AuthorID == me.ID:
		return t, httpx.Forbidden("You cannot review an item you wrote")
	case !t.Open:
		return t, httpx.Conflict("This item is " + string(t.Status) + " and no longer in the queue")
	}
	return t, nil
}

// ProviderQuality reports how reviewers judged each AI provider's
// candidates, per language (AI-5). Admins only.
func (s *Service) ProviderQuality(ctx context.Context) ([]dbgen.ProviderQuality, error) {
	me, err := accounts.RequireAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var rows []dbgen.ProviderQuality
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		rows, err = NewStore(tx).providerQuality(ctx)
		return err
	})
	return rows, err
}

// Credits lists the volunteers credited for accepted work in a language
// (CON-5), those who chose to be. Public.
func (s *Service) Credits(ctx context.Context, language string) ([]dbgen.ListCreditsRow, error) {
	var rows []dbgen.ListCreditsRow
	err := database.InTx(ctx, s.db, "", func(tx pgx.Tx) error {
		var err error
		rows, err = NewStore(tx).credits(ctx, language)
		return err
	})
	return rows, err
}
