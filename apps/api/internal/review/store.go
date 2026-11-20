package review

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

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

// decisions lists every decision on an item.
func (s *Store) decisions(ctx context.Context, itemID uuid.UUID) ([]dbgen.ReviewDecision, error) {
	out, err := s.q.ListItemDecisions(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("list decisions on %s: %w", itemID, err)
	}
	return out, nil
}
