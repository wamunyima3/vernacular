package review

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

// Dispute is a candidate reviewers disagree on, for a senior reviewer (CON-4).
type Dispute struct {
	Language      string
	Candidate     lexicon.SenseCard
	DisputedSince time.Time
	Reviews       []DisputeReview
}

// DisputeReview is one counted decision, without its reviewer.
type DisputeReview struct {
	Decision   dbgen.ReviewDecision
	Comment    string
	Correction string // the word the reviewer wrote instead
	CreatedAt  time.Time
}

// Disputes lists the disputes the signed-in senior reviewer can settle,
// oldest first.
func (s *Service) Disputes(ctx context.Context) ([]Dispute, error) {
	me, err := accounts.Require(ctx, nil)
	if err != nil {
		return nil, err
	}
	languages := me.SeniorLanguages()
	if len(languages) == 0 {
		return nil, httpx.ErrForbidden
	}
	var out []Dispute
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		store := NewStore(tx)
		rows, err := store.disputes(ctx, languages)
		if err != nil || len(rows) == 0 {
			return err
		}
		ids := make([]uuid.UUID, len(rows))
		for i, r := range rows {
			ids[i] = r.ID
		}
		reviews, err := store.counted(ctx, ids)
		if err != nil {
			return err
		}
		// Corrections name the candidate they created; show its word.
		wanted := slices.Clone(ids)
		for _, r := range reviews {
			if id := correctedTo(r.Correction); id != uuid.Nil {
				wanted = append(wanted, id)
			}
		}
		cards, err := newLexicon(tx).SenseCards(ctx, wanted)
		if err != nil {
			return err
		}
		byID := make(map[uuid.UUID]lexicon.SenseCard, len(cards))
		for _, c := range cards {
			byID[c.ID] = c
		}
		byItem := map[uuid.UUID][]DisputeReview{}
		for _, r := range reviews {
			byItem[r.ItemID] = append(byItem[r.ItemID], DisputeReview{
				Decision: r.Decision, Comment: r.Comment, CreatedAt: r.CreatedAt,
				Correction: byID[correctedTo(r.Correction)].Headword,
			})
		}
		out = make([]Dispute, len(rows))
		for i, r := range rows {
			out[i] = Dispute{Language: r.LanguageCode, Candidate: byID[r.ID], DisputedSince: r.DisputedSince, Reviews: byItem[r.ID]}
		}
		return nil
	})
	return out, err
}

// Resolve settles a dispute: `keep` verifies the item and `send_back`
// returns it to the queue. Seniors cannot settle what they wrote or reviewed.
func (s *Service) Resolve(ctx context.Context, senseID uuid.UUID, outcome dbgen.ResolutionOutcome, note string) (items.Status, error) {
	var status items.Status
	err := s.withNote(ctx, &note, func(tx pgx.Tx, me accounts.Principal) error {
		store, its := NewStore(tx), newItems(tx)
		item, ok, err := its.Lock(ctx, senseID)
		switch {
		case err != nil:
			return err
		case !ok || item.Kind != items.KindSense:
			return httpx.ErrNotFound
		case !me.IsSenior(item.Language):
			return httpx.ErrForbidden
		case item.Status != items.StatusDisputed:
			return httpx.Conflict("This item is " + string(item.Status) + ", not disputed")
		}
		reviewed, err := store.hasReviewed(ctx, senseID, me.ID)
		if err != nil {
			return err
		}
		if reviewed || item.AuthorID == me.ID {
			return httpx.Forbidden("You wrote or reviewed this item; another senior reviewer settles it")
		}
		if err := store.resolve(ctx, senseID, me.ID, outcome, note); err != nil {
			return err
		}
		status, err = promote(ctx, store, its, senseID, item.Origin)
		return err
	})
	return status, err
}

// withNote runs a signed-in reviewer's decision that must say why, in a
// transaction acting as them. note is trimmed in place.
func (s *Service) withNote(ctx context.Context, note *string, fn func(tx pgx.Tx, me accounts.Principal) error) error {
	me, err := accounts.Require(ctx, nil)
	if err != nil {
		return err
	}
	if *note = strings.TrimSpace(*note); *note == "" {
		return httpx.BadRequest("Write a note saying why")
	}
	return database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error { return fn(tx, me) })
}
