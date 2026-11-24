package review_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/phrases"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/review"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

// TestPhraseCards: phrase candidates get cards of their own, which take
// decisions and a reviewer's own phrase like word cards (P1-AI-04b).
func TestPhraseCards(t *testing.T) {
	tx := dbtest.Tx(t)
	phraseFixture(t, tx)
	svc := review.NewService(tx)
	card, ok, err := svc.Next(reviewer(t), "tst", nil)
	if err != nil || !ok || card.Kind != dbgen.ItemKindPhrase || len(card.Candidates) != 2 || card.GlossEn != "Zz good night" {
		t.Fatalf("next = %+v %v %v; want the phrase card with 2 phrases", card, ok, err)
	}

	checkPhraseCorrection(t, tx, svc, card)

	// Another reviewer sees the reviewer's phrase with the others.
	card, ok, err = svc.Next(reviewer(t, "r2"), "tst", nil)
	if err != nil || !ok || card.Kind != dbgen.ItemKindPhrase || len(card.Candidates) != 3 {
		t.Errorf("r2 next = %+v %v %v; want the 3 phrases", card, ok, err)
	}
}

// phraseFixture leaves one card in "tst": two Claude phrases for a line.
func phraseFixture(t *testing.T, tx pgx.Tx) {
	t.Helper()
	ctx := t.Context()
	fixture(t, tx)
	// Only the phrase card: the word cards are done.
	if _, err := tx.Exec(ctx, `UPDATE items SET status = 'verified' WHERE language_code = 'tst'`); err != nil {
		t.Fatal(err)
	}
	if _, err := phrases.ImportList(ctx, tx, strings.NewReader("pack,source_en,notes\nzz-greet,Zz good night,\n")); err != nil {
		t.Fatal(err)
	}
	store := phrases.NewStore(tx)
	var lineID uuid.UUID
	if err := tx.QueryRow(ctx, "SELECT id FROM phrase_list WHERE pack = 'zz-greet'").Scan(&lineID); err != nil {
		t.Fatal(err)
	}
	line, err := store.Line(ctx, lineID)
	if err != nil {
		t.Fatal(err)
	}
	claude, err := sources.Lookup(ctx, tx, "Claude")
	if err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"Zz mushileni", "Zz ulale bwino"} {
		if _, _, err := store.Propose(ctx, "tst", line, text, "", items.Machine(claude, "m", "p")); err != nil {
			t.Fatal(err)
		}
	}
}

// checkPhraseCorrection: writing a candidate as one's own is refused; a new
// phrase replaces the candidates, in the line's pack.
func checkPhraseCorrection(t *testing.T, tx pgx.Tx, svc *review.Service, card review.Card) {
	t.Helper()
	ctx := t.Context()
	all := func(d dbgen.ReviewDecision) []review.Decision {
		return []review.Decision{{ItemID: card.Candidates[0].ID, Decision: d}, {ItemID: card.Candidates[1].ID, Decision: d}}
	}
	if _, err := svc.Decide(reviewer(t), all(dbgen.ReviewDecisionCorrect), "", &review.Correction{Headword: "zz MUSHILENI"}); httpx.Status(err) != 400 {
		t.Errorf("correcting to a candidate: err = %v, want 400", err)
	}
	got, err := svc.Decide(reviewer(t), all(dbgen.ReviewDecisionCorrect), "", &review.Correction{Headword: "Zz mwalale"})
	if err != nil || len(got) != 2 || got[0].CorrectionItemID == nil {
		t.Fatalf("correct = %+v, %v", got, err)
	}
	var text, pack string
	if err := tx.QueryRow(ctx, "SELECT text, pack FROM phrases WHERE id = $1", *got[0].CorrectionItemID).Scan(&text, &pack); err != nil || text != "Zz mwalale" || pack != "zz-greet" {
		t.Errorf("correction = %q in %q, %v; want a new phrase in the line's pack", text, pack, err)
	}
}

// TestPhraseDisputes: a phrase reviewers disagree on reaches senior
// reviewers, who can settle it (P1-REV-04c).
func TestPhraseDisputes(t *testing.T) {
	tx := dbtest.Tx(t)
	phraseFixture(t, tx)
	svc := review.NewService(tx)
	card, ok, err := svc.Next(reviewer(t), "tst", nil)
	if err != nil || !ok {
		t.Fatalf("next = %v %v", ok, err)
	}
	phrase := card.Candidates[0].ID
	for _, r := range []struct {
		id string
		d  dbgen.ReviewDecision
	}{{"r1", dbgen.ReviewDecisionAccept}, {"r2", dbgen.ReviewDecisionReject}} {
		if _, err := svc.Decide(reviewer(t, r.id), []review.Decision{{ItemID: phrase, Decision: r.d}}, "zz", nil); err != nil {
			t.Fatal(err)
		}
	}
	disputes, err := svc.Disputes(senior(t, "s1"))
	if err != nil || len(disputes) != 1 || disputes[0].Candidate.ID != phrase || disputes[0].Candidate.English != "Zz good night" {
		t.Fatalf("disputes = %+v, %v; want the phrase", disputes, err)
	}
	if status, err := svc.Resolve(senior(t, "s1"), phrase, dbgen.ResolutionOutcomeKeep, "Zz right"); err != nil || status != items.StatusVerified {
		t.Errorf("resolve = %s, %v; want verified", status, err)
	}
}
