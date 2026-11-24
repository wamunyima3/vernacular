package review_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/contributions"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/review"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

// TestCorrectionRetires: reviewers see what a correction would replace, and
// verifying it retires the original, which keeps its row (P1-REV-01b2).
func TestCorrectionRetires(t *testing.T) {
	tx := dbtest.Tx(t)
	ctx := t.Context()
	water, _ := fixture(t, tx)
	if _, err := tx.Exec(ctx, `UPDATE items SET status = 'verified' WHERE language_code = 'tst'`); err != nil {
		t.Fatal(err)
	}
	fix, err := contributions.Add(ctx, tx, "r3", contributions.Input{
		Language: "tst", Headword: "amensi", PartOfSpeech: "noun", GlossEn: "zz water",
		Corrects: uuid.NullUUID{UUID: water, Valid: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	svc := review.NewService(tx)
	card, ok, err := svc.Next(reviewer(t), "tst", nil)
	if err != nil || !ok || len(card.Candidates) != 1 {
		t.Fatalf("next = %+v %v %v; want the correction", card, ok, err)
	}
	if c := card.Candidates[0].Corrects; c == nil || c.ID != water || c.Headword != "amenshi" {
		t.Fatalf("corrects = %+v, want amenshi", c)
	}
	retired := func() (r bool) {
		if err := tx.QueryRow(ctx, `SELECT retired($1)`, water).Scan(&r); err != nil {
			t.Fatal(err)
		}
		return r
	}
	for _, id := range []string{"r1", "r2"} {
		if retired() {
			t.Fatalf("retired before %s accepted", id)
		}
		if _, err := svc.Decide(reviewer(t, id), []review.Decision{{ItemID: fix.SenseID, Decision: dbgen.ReviewDecisionAccept}}, "", nil); err != nil {
			t.Fatal(err)
		}
	}
	if !retired() {
		t.Error("the original is still served after its correction was verified")
	}
}
