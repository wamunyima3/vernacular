package accounts

import (
	"context"
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

// requireSenior admits the language's senior reviewers, and admins, so a
// language whose seniors are all suspended can still recover.
func requireSenior(ctx context.Context, language string) (Principal, error) {
	return Require(ctx, func(p Principal) bool { return p.IsSenior(language) || p.Role == gen.RoleAdmin })
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
