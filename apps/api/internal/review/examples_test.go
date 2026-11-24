package review_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/review"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

// TestExampleCards: examples get cards of their own showing the word they
// use; two accepts verify one, and a correction adds a new example to the
// same meaning (LEX-5, P2-09a).
func TestExampleCards(t *testing.T) {
	tx := dbtest.Tx(t)
	water := exampleFixture(t, tx)
	svc := review.NewService(tx)
	card, ok, err := svc.Next(reviewer(t), "tst", nil)
	if err != nil || !ok || card.Kind != dbgen.ItemKindExample || len(card.Candidates) != 2 || card.GlossEn != "zz water" {
		t.Fatalf("next = %+v %v %v; want the example card with 2 examples", card, ok, err)
	}
	if c := card.Candidates[0]; c.Word != "amenshi" || c.TranslationEn == "" {
		t.Errorf("candidate = %+v; want the word and the English", c)
	}

	first := card.Candidates[0].ID
	for _, who := range []string{"r1", "r2"} {
		if _, err := svc.Decide(reviewer(t, who), []review.Decision{{ItemID: first, Decision: dbgen.ReviewDecisionAccept}}, "", nil); err != nil {
			t.Fatal(err)
		}
	}
	var status string
	if err := tx.QueryRow(t.Context(), "SELECT status FROM items WHERE id = $1", first).Scan(&status); err != nil || status != "verified" {
		t.Errorf("after two accepts: %s, %v; want verified", status, err)
	}

	checkExampleCorrection(t, tx, svc, card.Candidates[1], water)
}

// exampleFixture leaves one card in "tst": two contributed examples of the
// water sense, by r3 so r1 and r2 may review them.
func exampleFixture(t *testing.T, tx pgx.Tx) uuid.UUID {
	t.Helper()
	ctx := t.Context()
	water, _ := fixture(t, tx)
	if _, err := tx.Exec(ctx, `UPDATE items SET status = 'verified' WHERE language_code = 'tst'`); err != nil {
		t.Fatal(err)
	}
	people, err := sources.Lookup(ctx, tx, "Community contributors")
	if err != nil {
		t.Fatal(err)
	}
	lex := lexicon.NewStore(tx)
	for _, x := range [][2]string{{"Zz amenshi yali ayasuma", "The water was good"}, {"Zz nalanwa amenshi", "I drank water"}} {
		if _, err := lex.AddExample(ctx, "tst", water, x[0], x[1], items.Contributed(people, "r3")); err != nil {
			t.Fatal(err)
		}
	}
	return water
}

// checkExampleCorrection: writing the candidate as one's own is refused; a new
// sentence becomes an example of the same meaning.
func checkExampleCorrection(t *testing.T, tx pgx.Tx, svc *review.Service, second review.Candidate, water uuid.UUID) {
	t.Helper()
	if _, err := svc.Decide(reviewer(t), []review.Decision{{ItemID: second.ID, Decision: dbgen.ReviewDecisionCorrect}}, "",
		&review.Correction{Headword: second.Text, TranslationEn: "x"}); httpx.Status(err) != 400 {
		t.Errorf("correcting to the candidate: %v, want 400", err)
	}
	got, err := svc.Decide(reviewer(t), []review.Decision{{ItemID: second.ID, Decision: dbgen.ReviewDecisionCorrect}}, "",
		&review.Correction{Headword: "Zz amenshi yalitalala", TranslationEn: "The water is cold"})
	if err != nil || got[0].CorrectionItemID == nil {
		t.Fatalf("correct = %+v, %v", got, err)
	}
	var sense uuid.UUID
	if err := tx.QueryRow(t.Context(), "SELECT sense_id FROM examples WHERE id = $1", *got[0].CorrectionItemID).Scan(&sense); err != nil || sense != water {
		t.Errorf("correction on sense %v, %v; want the water sense", sense, err)
	}
}
