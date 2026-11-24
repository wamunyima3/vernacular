package review

import (
	"context"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

// Status is the promotion rule (ADR-0004, DQ-2): an item's trust status from
// its origin, a senior reviewer's latest resolution ("" when none), whether a
// confirmed report awaits a resolution, and the decisions that count, those
// no send-back set aside. It is the only code that decides `verified`.
// Reviews are already distinct, approved and never the author's, because
// Decide refuses anything else.
func Status(origin items.Status, resolution dbgen.ResolutionOutcome, reported bool, decisions []dbgen.ReviewDecision) items.Status {
	switch {
	case reported:
		return items.StatusDisputed
	case resolution == dbgen.ResolutionOutcomeKeep:
		return items.StatusVerified
	}
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

// sendBack returns a locked item to the review queue, as a senior's
// send-back resolution does: reviews so far are set aside.
func sendBack(ctx context.Context, store *Store, its Items, itemID uuid.UUID, origin items.Status, actorID, note string) (items.Status, error) {
	if err := store.resolve(ctx, itemID, actorID, dbgen.ResolutionOutcomeSendBack, note); err != nil {
		return "", err
	}
	return promote(ctx, store, its, itemID, origin)
}

// promote recomputes an item's status after its reviews or resolutions
// change. The caller holds the item's lock, which is where origin comes from.
func promote(ctx context.Context, store *Store, its Items, itemID uuid.UUID, origin items.Status) (items.Status, error) {
	resolution, reported, decisions, err := store.promotionInputs(ctx, itemID)
	if err != nil {
		return "", err
	}
	status := Status(origin, resolution, reported, decisions)
	return status, its.SetStatus(ctx, itemID, status)
}
