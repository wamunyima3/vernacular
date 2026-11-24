package review

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

// Quarter names the calendar quarter t falls in, in UTC, as "2026-Q4".
func Quarter(t time.Time) string {
	start := quarterStart(t)
	return fmt.Sprintf("%d-Q%d", start.Year(), int(start.Month())/3+1)
}

// quarterStart is the first moment of t's calendar quarter, UTC.
func quarterStart(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), (t.Month()-1)/3*3+1, 1, 0, 0, 0, 0, time.UTC)
}

// Sample draws the quarter's re-check sample (DQ-6, US-28): 5% of each
// language's verified items, at least one. A quarter is drawn once; later
// calls return 0. An admin starts it by hand, the worker each quarter.
func (s *Service) Sample(ctx context.Context, quarter string) (int64, error) {
	me, err := accounts.RequireAdmin(ctx)
	if err != nil {
		return 0, err
	}
	return s.sample(ctx, me.ID, quarter)
}

// SampleAsSystem is Sample for the worker.
func (s *Service) SampleAsSystem(ctx context.Context, quarter string) (int64, error) {
	return s.sample(ctx, "system:recheck", quarter)
}

func (s *Service) sample(ctx context.Context, actor, quarter string) (int64, error) {
	var n int64
	err := database.InTx(ctx, s.db, actor, func(tx pgx.Tx) error {
		var err error
		n, err = dbgen.New(tx).SampleRechecks(ctx, quarter)
		return err
	})
	return n, err
}

// Recheck is one sampled item waiting to be checked again.
type Recheck struct {
	Quarter   string
	Candidate Candidate
}

// NextRecheck returns the oldest open re-check in a language for a senior
// reviewer who neither wrote nor reviewed the item, or false when none is left.
func (s *Service) NextRecheck(ctx context.Context, language string, skip []uuid.UUID) (Recheck, bool, error) {
	me, err := accounts.Require(ctx, func(p accounts.Principal) bool { return p.IsSenior(language) })
	if err != nil {
		return Recheck{}, false, err
	}
	var out Recheck
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		r, err := dbgen.New(tx).NextRecheck(ctx, dbgen.NextRecheckParams{Language: language, ReviewerID: me.ID, Skip: skip})
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("next re-check: %w", err)
		}
		cands, err := loadCandidates(ctx, tx, []uuid.UUID{r.ItemID})
		if err != nil || len(cands) == 0 {
			return err
		}
		out = Recheck{Quarter: r.Quarter, Candidate: cands[0]}
		return nil
	})
	return out, out.Quarter != "", err
}

// Checked is a recorded re-check and the item's status after it.
type Checked struct {
	Quarter string
	Status  items.Status
}

// RecordRecheck records a senior reviewer's re-check. A pass changes
// nothing; a fail, which says why, sends the item back to the review queue
// as a send-back resolution does, so only reviews made after it count.
func (s *Service) RecordRecheck(ctx context.Context, itemID uuid.UUID, outcome dbgen.RecheckOutcome, note string) (Checked, error) {
	me, err := accounts.Require(ctx, nil)
	if err != nil {
		return Checked{}, err
	}
	note = strings.TrimSpace(note)
	if outcome == dbgen.RecheckOutcomeFail && note == "" {
		return Checked{}, httpx.BadRequest("Write a note saying what is wrong")
	}
	var out Checked
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		q, store, its := dbgen.New(tx), NewStore(tx), newItems(tx)
		item, ok, err := its.Lock(ctx, itemID)
		if err != nil {
			return err
		}
		open, err := q.LockOpenRecheck(ctx, itemID)
		switch {
		case !ok || errors.Is(err, pgx.ErrNoRows):
			return httpx.ErrNotFound
		case err != nil:
			return fmt.Errorf("lock re-check of %s: %w", itemID, err)
		case !me.IsSenior(open.LanguageCode):
			return httpx.ErrForbidden
		}
		if err := store.notOwn(ctx, itemID, item.AuthorID, me.ID, "re-checks"); err != nil {
			return err
		}
		err = q.RecordRecheck(ctx, dbgen.RecordRecheckParams{
			Quarter: open.Quarter, ItemID: itemID, Outcome: dbgen.NullRecheckOutcome{RecheckOutcome: outcome, Valid: true},
			Note: database.Text(note), CheckedBy: database.Text(me.ID),
		})
		if err != nil {
			return fmt.Errorf("record re-check of %s: %w", itemID, err)
		}
		out = Checked{Quarter: open.Quarter, Status: item.Status}
		if outcome == dbgen.RecheckOutcomePass {
			return nil
		}
		out.Status, err = sendBack(ctx, store, its, itemID, item.Origin, me.ID, "Re-check "+open.Quarter+": "+note)
		return err
	})
	return out, err
}
