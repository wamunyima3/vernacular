package review_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/review"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

const (
	accept  = dbgen.ReviewDecisionAccept
	reject  = dbgen.ReviewDecisionReject
	correct = dbgen.ReviewDecisionCorrect
)

func TestStatus(t *testing.T) {
	tests := []struct {
		name      string
		origin    items.Status
		decisions []dbgen.ReviewDecision
		want      items.Status
	}{
		{"unreviewed machine", items.StatusMachine, nil, items.StatusMachine},
		{"unreviewed sourced", items.StatusSourced, nil, items.StatusSourced},
		{"one accept", items.StatusMachine, []dbgen.ReviewDecision{accept}, items.StatusReviewed},
		{"two accepts", items.StatusMachine, []dbgen.ReviewDecision{accept, accept}, items.StatusVerified},
		{"three accepts", items.StatusSourced, []dbgen.ReviewDecision{accept, accept, accept}, items.StatusVerified},
		{"one reject waits for another opinion", items.StatusSourced, []dbgen.ReviewDecision{reject}, items.StatusSourced},
		{"two rejects", items.StatusMachine, []dbgen.ReviewDecision{reject, reject}, items.StatusRejected},
		{"a correction counts as a reject", items.StatusMachine, []dbgen.ReviewDecision{reject, correct}, items.StatusRejected},
		{"accept then reject", items.StatusMachine, []dbgen.ReviewDecision{accept, reject}, items.StatusDisputed},
		{"two accepts and a reject", items.StatusSourced, []dbgen.ReviewDecision{accept, accept, reject}, items.StatusDisputed},
		{"accept and correction", items.StatusSourced, []dbgen.ReviewDecision{correct, accept}, items.StatusDisputed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := review.Status(tt.origin, tt.decisions); got != tt.want {
				t.Errorf("Status(%s, %v) = %s, want %s", tt.origin, tt.decisions, got, tt.want)
			}
		})
	}
}

// TestPromotion runs the full cycle: two independent accepts verify an item,
// undoing the second takes it back to reviewed, and a disagreement disputes it.
func TestPromotion(t *testing.T) {
	tx := dbtest.Tx(t)
	water, _ := fixture(t, tx)
	if _, err := tx.Exec(t.Context(), `INSERT INTO users (id, name, email, email_verified, role)
		VALUES ('r2', 'Second', 'r2@test.invalid', true, 'contributor'), ('r3', 'Third', 'r3@test.invalid', true, 'contributor')`); err != nil {
		t.Fatal(err)
	}
	svc := review.NewService(tx)
	status := func() items.Status {
		item, _, err := items.NewStore(tx).Lock(t.Context(), water)
		if err != nil {
			t.Fatal(err)
		}
		return item.Status
	}
	decide := func(who string, d dbgen.ReviewDecision) []uuid.UUID {
		t.Helper()
		r, err := svc.Decide(reviewer(t, who), []review.Decision{{SenseID: water, Decision: d}}, "", nil)
		if err != nil {
			t.Fatalf("%s %s: %v", who, d, err)
		}
		return []uuid.UUID{r[0].ID}
	}

	decide("r1", accept)
	if got := status(); got != items.StatusReviewed {
		t.Fatalf("after one accept: %s, want reviewed", got)
	}
	second := decide("r2", accept)
	if got := status(); got != items.StatusVerified {
		t.Fatalf("after two accepts: %s, want verified", got)
	}
	if card, ok, _ := svc.Next(reviewer(t, "r3"), "tst", nil); ok && card.GlossEn == "zz water" {
		t.Error("a verified item is still in the queue")
	}
	if err := svc.Undo(reviewer(t, "r2"), second); err != nil {
		t.Fatal(err)
	}
	if got := status(); got != items.StatusReviewed {
		t.Fatalf("after undoing the second accept: %s, want reviewed", got)
	}
	decide("r3", reject)
	if got := status(); got != items.StatusDisputed {
		t.Fatalf("after accept and reject: %s, want disputed", got)
	}
}
