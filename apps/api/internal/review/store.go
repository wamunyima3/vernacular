package review

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

// errAlreadyReviewed means the reviewer has decided on the item before.
var errAlreadyReviewed = errors.New("already reviewed")

// Store reads the queue and writes reviews.
type Store struct {
	q *dbgen.Queries
}

// NewStore wraps a pool or transaction.
func NewStore(db dbgen.DBTX) *Store { return &Store{q: dbgen.New(db)} }

// nextCard returns the senses of the next card, in display order. Empty
// when the queue is.
func (s *Store) nextCard(ctx context.Context, language, reviewerID string, skip []uuid.UUID) ([]dbgen.NextReviewCardRow, error) {
	rows, err := s.q.NextReviewCard(ctx, dbgen.NextReviewCardParams{LanguageCode: language, ReviewerID: reviewerID, Skip: skip})
	if err != nil {
		return nil, fmt.Errorf("next review card: %w", err)
	}
	return rows, nil
}

// create records a decision; correction is the JSON detail of a `correct` one.
func (s *Store) create(ctx context.Context, itemID uuid.UUID, reviewerID string, decision dbgen.ReviewDecision, comment string, correction []byte) (uuid.UUID, error) {
	id, err := s.q.CreateReview(ctx, dbgen.CreateReviewParams{
		ItemID: itemID, ReviewerID: reviewerID, Decision: decision, Comment: comment, Correction: correction,
	})
	if database.Violates(err, database.UniqueViolation) {
		return uuid.Nil, errAlreadyReviewed
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("create review of %s: %w", itemID, err)
	}
	return id, nil
}

// own returns the reviewer's reviews among ids, with what they are about.
func (s *Store) own(ctx context.Context, ids []uuid.UUID, reviewerID string) ([]dbgen.ListOwnReviewsRow, error) {
	rows, err := s.q.ListOwnReviews(ctx, dbgen.ListOwnReviewsParams{Ids: ids, ReviewerID: reviewerID})
	if err != nil {
		return nil, fmt.Errorf("list own reviews: %w", err)
	}
	return rows, nil
}

func (s *Store) delete(ctx context.Context, ids []uuid.UUID) error {
	if err := s.q.DeleteReviews(ctx, ids); err != nil {
		return fmt.Errorf("delete reviews: %w", err)
	}
	return nil
}

// promotionInputs reads what the promotion rule needs besides origin: the
// latest resolution, if any, and the decisions that count.
func (s *Store) promotionInputs(ctx context.Context, itemID uuid.UUID) (dbgen.ResolutionOutcome, bool, []dbgen.ReviewDecision, error) {
	var resolution dbgen.ResolutionOutcome
	latest, err := s.q.GetLatestResolution(ctx, itemID)
	switch {
	case err == nil:
		resolution = latest
	case !errors.Is(err, pgx.ErrNoRows):
		return "", false, nil, fmt.Errorf("latest resolution of %s: %w", itemID, err)
	}
	reported, err := s.q.HasUnsettledReport(ctx, itemID)
	if err != nil {
		return "", false, nil, fmt.Errorf("unsettled reports of %s: %w", itemID, err)
	}
	reviews, err := s.counted(ctx, []uuid.UUID{itemID})
	if err != nil {
		return "", false, nil, err
	}
	decisions := make([]dbgen.ReviewDecision, len(reviews))
	for i, r := range reviews {
		decisions[i] = r.Decision
	}
	return resolution, reported, decisions, nil
}

// counted lists the reviews of items that count, oldest first per item.
func (s *Store) counted(ctx context.Context, ids []uuid.UUID) ([]dbgen.ListCountedReviewsRow, error) {
	rows, err := s.q.ListCountedReviews(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("list counted reviews: %w", err)
	}
	return rows, nil
}

func (s *Store) disputes(ctx context.Context, languages []string) ([]dbgen.ListDisputesRow, error) {
	rows, err := s.q.ListDisputes(ctx, languages)
	if err != nil {
		return nil, fmt.Errorf("list disputes: %w", err)
	}
	return rows, nil
}

// resolve records a resolution; a send-back sets the item's reviews aside.
func (s *Store) resolve(ctx context.Context, itemID uuid.UUID, seniorID string, outcome dbgen.ResolutionOutcome, note string) error {
	id, err := s.q.CreateResolution(ctx, dbgen.CreateResolutionParams{ItemID: itemID, SeniorID: seniorID, Outcome: outcome, Note: note})
	if err != nil {
		return fmt.Errorf("resolve %s: %w", itemID, err)
	}
	resolution := uuid.NullUUID{UUID: id, Valid: true}
	if err := s.q.SettleReports(ctx, dbgen.SettleReportsParams{ResolutionID: resolution, ItemID: itemID}); err != nil {
		return fmt.Errorf("settle reports of %s: %w", itemID, err)
	}
	if outcome == dbgen.ResolutionOutcomeSendBack {
		err = s.q.SetReviewsAside(ctx, dbgen.SetReviewsAsideParams{ResolutionID: resolution, ItemID: itemID})
		if err != nil {
			return fmt.Errorf("set reviews of %s aside: %w", itemID, err)
		}
	}
	return nil
}

func (s *Store) hasReviewed(ctx context.Context, itemID uuid.UUID, reviewerID string) (bool, error) {
	ok, err := s.q.HasReviewed(ctx, dbgen.HasReviewedParams{ItemID: itemID, ReviewerID: reviewerID})
	if err != nil {
		return false, fmt.Errorf("has %s reviewed %s: %w", reviewerID, itemID, err)
	}
	return ok, nil
}

func (s *Store) providerQuality(ctx context.Context) ([]dbgen.ProviderQuality, error) {
	rows, err := s.q.ListProviderQuality(ctx)
	if err != nil {
		return nil, fmt.Errorf("list provider quality: %w", err)
	}
	return rows, nil
}

func (s *Store) credits(ctx context.Context, language string) ([]dbgen.ListCreditsRow, error) {
	rows, err := s.q.ListCredits(ctx, language)
	if err != nil {
		return nil, fmt.Errorf("list credits: %w", err)
	}
	return rows, nil
}
