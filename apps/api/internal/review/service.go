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
	Language   string
	ConceptID  uuid.UUID
	GlossEn    string // the concept's gloss, or on a phrase card the phrase's English line
	Priority   dbgen.ReviewPriority
	Kind       dbgen.ItemKind // a card holds one kind: senses, phrases or examples
	Candidates []Candidate
}

// Candidate is one item reviewers judge: a sense, a phrase or an example.
type Candidate struct {
	ID                                     uuid.UUID
	Kind                                   dbgen.ItemKind
	ConceptID                              uuid.NullUUID
	English                                string // the concept's gloss, or a phrase's English line
	Text                                   string // the headword, the phrase, or the example sentence
	Word, TranslationEn                    string // for an example: the word it uses, and its English
	PartOfSpeech, NounClass, Plural, Notes string
	Corrects                               *lexicon.Corrected // the meaning a correction would replace
}

// loadCandidates reads senses, phrases and examples by ID, in the order of ids;
// unknown IDs are left out.
func loadCandidates(ctx context.Context, tx pgx.Tx, ids []uuid.UUID) ([]Candidate, error) {
	senses, err := newLexicon(tx).SenseCards(ctx, ids)
	if err != nil {
		return nil, err
	}
	var (
		phr []phrases.Phrase
		exs []lexicon.ExampleCard
	)
	if len(senses) < len(ids) { // the rest may be phrases or examples
		if phr, err = newPhrases(tx).Cards(ctx, ids); err != nil {
			return nil, err
		}
	}
	if len(senses)+len(phr) < len(ids) {
		if exs, err = newLexicon(tx).ExampleCards(ctx, ids); err != nil {
			return nil, err
		}
	}
	byID := make(map[uuid.UUID]Candidate, len(ids))
	for _, c := range senses {
		byID[c.ID] = Candidate{
			ID: c.ID, Kind: dbgen.ItemKindSense, ConceptID: c.ConceptID, English: c.GlossEn, Text: c.Headword,
			PartOfSpeech: c.PartOfSpeech, NounClass: c.NounClass, Plural: c.Plural, Corrects: c.Corrects,
		}
	}
	for _, p := range phr {
		byID[p.ID] = Candidate{ID: p.ID, Kind: dbgen.ItemKindPhrase, ConceptID: p.ConceptID, English: p.SourceEn, Text: p.Text, Notes: p.Notes}
	}
	for _, x := range exs {
		byID[x.ID] = Candidate{
			ID: x.ID, Kind: dbgen.ItemKindExample, ConceptID: x.ConceptID, English: x.GlossEn, Text: x.Text,
			Word: x.Headword, TranslationEn: x.TranslationEn,
		}
	}
	out := make([]Candidate, 0, len(ids))
	for _, id := range ids {
		if c, ok := byID[id]; ok {
			out = append(out, c)
		}
	}
	return out, nil
}

// Decision is a reviewer's verdict on one candidate.
type Decision struct {
	ItemID   uuid.UUID
	Decision dbgen.ReviewDecision
}

// Recorded is a stored decision.
type Recorded struct {
	ID uuid.UUID
	Decision
	CorrectionItemID *uuid.UUID // the new candidate, for corrections
}

// Lexicon is what the queue needs from the lexicon module.
type Lexicon interface {
	SenseCards(ctx context.Context, ids []uuid.UUID) ([]lexicon.SenseCard, error)
	ExampleCards(ctx context.Context, ids []uuid.UUID) ([]lexicon.ExampleCard, error)
	AddExample(ctx context.Context, language string, senseID uuid.UUID, text, translationEn string, origin items.Origin) (uuid.UUID, error)
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
		cands, err := loadCandidates(ctx, tx, ids)
		if err != nil || len(cands) == 0 {
			return err
		}
		card = Card{
			Language: language, ConceptID: rows[0].ConceptID, GlossEn: cands[0].English,
			Priority: rows[0].Priority, Kind: rows[0].Kind, Candidates: cands,
		}
		return nil
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
	sorted := slices.SortedFunc(slices.Values(decisions), byItem)
	out := make([]Recorded, 0, len(decisions))
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		store, its := NewStore(tx), newItems(tx)
		locked := make([]items.Locked, len(sorted))
		for i, d := range sorted {
			if locked[i], err = checkTarget(ctx, its, me, d.ItemID); err != nil {
				return err
			}
			if locked[i].Kind != locked[0].Kind {
				return httpx.BadRequest("Decide on one card's candidates together: a card holds one kind of item")
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
				r.CorrectionItemID, detail = &fix.itemID, fix.detail
			}
			if r.ID, err = store.create(ctx, d.ItemID, me.ID, d.Decision, comment, detail); errors.Is(err, errAlreadyReviewed) {
				return httpx.Conflict("You have already reviewed " + d.ItemID.String())
			} else if err != nil {
				return err
			}
			if _, err := promote(ctx, store, its, d.ItemID, locked[i].Origin); err != nil {
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

func byItem(a, b Decision) int { return compareIDs(a.ItemID, b.ItemID) }

// correct records the reviewer's own word or phrase when the decisions
// include corrections; nil when they do not.
func correct(ctx context.Context, tx pgx.Tx, reviewerID string, locked []items.Locked, sorted []Decision, c *Correction) (*replacement, error) {
	if c == nil {
		return nil, nil
	}
	switch locked[0].Kind {
	case items.KindPhrase:
		return replacePhrase(ctx, tx, reviewerID, locked[0].Language, sorted, *c)
	case items.KindExample:
		return replaceExample(ctx, tx, reviewerID, locked[0].Language, sorted, *c)
	}
	return replace(ctx, tx, reviewerID, locked[0].Language, sorted, *c)
}

// checkTarget locks one item and applies the blind-review rules (CON-2, CON-3).
func checkTarget(ctx context.Context, its Items, me accounts.Principal, id uuid.UUID) (items.Locked, error) {
	t, ok, err := its.Lock(ctx, id)
	switch {
	case err != nil:
		return t, err
	case !ok || !items.Reviewable(t.Kind):
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
