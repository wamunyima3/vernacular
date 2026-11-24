package accounts

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// Service lets senior reviewers monitor the reviewers of their languages
// (DQ-7, US-24).
type Service struct {
	db database.Beginner
}

// NewService returns an accounts service.
func NewService(db database.Beginner) *Service { return &Service{db: db} }

// requireSenior admits those who may act as a senior reviewer of the language.
func requireSenior(ctx context.Context, language string) (Principal, error) {
	return Require(ctx, func(p Principal) bool { return p.SeniorOrAdmin(language) })
}

// Reviewers lists a language's reviewers with their agreement, suspended first.
func (s *Service) Reviewers(ctx context.Context, language string) ([]dbgen.ListReviewersRow, error) {
	me, err := requireSenior(ctx, language)
	if err != nil {
		return nil, err
	}
	var rows []dbgen.ListReviewersRow
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		rows, err = dbgen.New(tx).ListReviewers(ctx, language)
		if err != nil {
			return fmt.Errorf("list reviewers: %w", err)
		}
		return nil
	})
	return rows, err
}

// Retrain records that a reviewer was retrained, which restarts their
// agreement count and lifts a suspension. Nobody retrains themselves.
func (s *Service) Retrain(ctx context.Context, userID, language string) error {
	me, err := requireSenior(ctx, language)
	if err != nil {
		return err
	}
	if userID == me.ID {
		return httpx.Forbidden("Another senior reviewer records your retraining")
	}
	return database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		n, err := dbgen.New(tx).MarkRetrained(ctx, dbgen.MarkRetrainedParams{UserID: userID, LanguageCode: language})
		if err != nil {
			return fmt.Errorf("mark %s retrained: %w", userID, err)
		}
		if n == 0 {
			return httpx.ErrNotFound
		}
		return nil
	})
}

// Approvals lists every reviewer approval, for the admin (CON-2).
func (s *Service) Approvals(ctx context.Context) ([]dbgen.ListApprovalsRow, error) {
	me, err := RequireAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var rows []dbgen.ListApprovalsRow
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		if rows, err = dbgen.New(tx).ListApprovals(ctx); err != nil {
			return fmt.Errorf("list approvals: %w", err)
		}
		return nil
	})
	return rows, err
}

// Approve lets the user with the email review a language at a level, or
// changes their level.
func (s *Service) Approve(ctx context.Context, email, language string, level gen.ReviewerLevel) error {
	me, err := RequireAdmin(ctx)
	if err != nil {
		return err
	}
	return database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		_, err := dbgen.New(tx).UpsertApproval(ctx, dbgen.UpsertApprovalParams{
			LanguageCode: language, Level: dbgen.ReviewerLevel(level), ApprovedBy: database.Text(me.ID), Email: email,
		})
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return httpx.BadRequest("No user has signed up with " + email)
		case database.Violates(err, database.ForeignKeyViolation):
			return httpx.BadRequest("No language " + language)
		case err != nil:
			return fmt.Errorf("approve %s: %w", email, err)
		}
		return nil
	})
}

// Revoke removes a user's approval for a language.
func (s *Service) Revoke(ctx context.Context, userID, language string) error {
	me, err := RequireAdmin(ctx)
	if err != nil {
		return err
	}
	return database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		n, err := dbgen.New(tx).DeleteApproval(ctx, dbgen.DeleteApprovalParams{UserID: userID, LanguageCode: language})
		if err != nil {
			return fmt.Errorf("revoke %s for %s: %w", userID, language, err)
		}
		if n == 0 {
			return httpx.ErrNotFound
		}
		return nil
	})
}

// SetCredit records whether the signed-in user wants public credit by name.
func (s *Service) SetCredit(ctx context.Context, credited bool) error {
	me, err := Require(ctx, nil)
	if err != nil {
		return err
	}
	return database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		if err := dbgen.New(tx).SetCreditConsent(ctx, dbgen.SetCreditConsentParams{Credited: credited, ID: me.ID}); err != nil {
			return fmt.Errorf("set credit consent: %w", err)
		}
		return nil
	})
}
