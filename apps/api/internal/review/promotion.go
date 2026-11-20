package review

import (
	"context"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

// Status is the promotion rule (ADR-0004, DQ-2): an item's trust status from
// its origin and the decisions of its reviewers. It is the only code that
// decides `verified`. Reviews are already distinct, approved and never the
// author's, because Decide refuses anything else.
func Status(origin items.Status, decisions []dbgen.ReviewDecision) items.Status {
	var accepts, rejects int
	for _, d := range decisions {
		if d == dbgen.ReviewDecisionAccept {
			accepts++
		} else {
			rejects++ // a correction says the item was wrong, like a reject
		}
	}
	switch {
	case accepts > 0 && rejects > 0:
		return items.StatusDisputed
	case accepts >= 2:
		return items.StatusVerified
	case accepts == 1:
		return items.StatusReviewed
	case rejects >= 2:
		return items.StatusRejected
	default:
		return origin
	}
}

// promote recomputes an item's status after its reviews change. The caller
// holds the item's lock, which is where origin comes from.
func promote(ctx context.Context, store *Store, its Items, itemID uuid.UUID, origin items.Status) error {
	decisions, err := store.decisions(ctx, itemID)
	if err != nil {
		return err
	}
	return its.SetStatus(ctx, itemID, Status(origin, decisions))
}
