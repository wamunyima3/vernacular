package review

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

// Report records that someone thinks a sense is wrong (PRV-3, US-15). The
// reporter may be signed out.
func (s *Service) Report(ctx context.Context, senseID uuid.UUID, reason string) (uuid.UUID, error) {
	me, _ := accounts.FromContext(ctx)
	if me.ByKey() { // a key names an app, such as the portal, not a person
		me = accounts.Principal{}
	}
	if reason = strings.TrimSpace(reason); reason == "" {
		return uuid.Nil, httpx.BadRequest("Say what is wrong")
	}
	var id uuid.UUID
	err := database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		var err error
		id, err = dbgen.New(tx).CreateReport(ctx, dbgen.CreateReportParams{ItemID: senseID, Reason: reason, ReporterID: database.Text(me.ID)})
		if errors.Is(err, pgx.ErrNoRows) {
			return httpx.ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("create report: %w", err)
		}
		return nil
	})
	return id, err
}

// Reports lists a language's open reports, for its reviewers.
func (s *Service) Reports(ctx context.Context, language string) ([]dbgen.ListOpenReportsRow, error) {
	me, err := accounts.Require(ctx, func(p accounts.Principal) bool { return p.CanReview(language) })
	if err != nil {
		return nil, err
	}
	var rows []dbgen.ListOpenReportsRow
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		rows, err = dbgen.New(tx).ListOpenReports(ctx, []string{language})
		if err != nil {
			return fmt.Errorf("list reports: %w", err)
		}
		return nil
	})
	return rows, err
}

// DecideReport confirms or dismisses an open report. Confirming disputes the
// sense until a senior reviewer settles it (US-15). It returns the sense's
// status afterwards.
func (s *Service) DecideReport(ctx context.Context, reportID uuid.UUID, confirm bool, note string) (items.Status, error) {
	var status items.Status
	err := s.withNote(ctx, &note, func(tx pgx.Tx, me accounts.Principal) error {
		q, its := dbgen.New(tx), newItems(tx)
		itemID, err := q.GetReportItem(ctx, reportID)
		if errors.Is(err, pgx.ErrNoRows) {
			return httpx.ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("get report: %w", err)
		}
		// The item first, then the report: Resolve locks the item before it
		// settles reports, so the same order cannot deadlock.
		item, ok, err := its.Lock(ctx, itemID)
		switch {
		case err != nil:
			return err
		case !ok:
			return httpx.ErrNotFound
		case !me.CanReview(item.Language):
			return httpx.ErrForbidden
		}
		report, err := q.LockReport(ctx, reportID)
		switch {
		case err != nil:
			return fmt.Errorf("lock report: %w", err)
		case report.State != dbgen.ReportStateOpen:
			return httpx.Conflict("This report was already " + string(report.State))
		case report.ReporterID.String == me.ID:
			return httpx.Forbidden("Another reviewer decides your own reports")
		}
		state := dbgen.ReportStateDismissed
		if confirm {
			state = dbgen.ReportStateConfirmed
		}
		if err := q.HandleReport(ctx, dbgen.HandleReportParams{ID: reportID, State: state, HandledBy: database.Text(me.ID), Note: note}); err != nil {
			return fmt.Errorf("handle report: %w", err)
		}
		status, err = promote(ctx, NewStore(tx), its, itemID, item.Origin)
		return err
	})
	return status, err
}
