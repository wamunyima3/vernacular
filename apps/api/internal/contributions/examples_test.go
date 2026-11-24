package contributions_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/contributions"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// TestExampleCoverage: a contributor adds an example to an existing meaning,
// and a verified core meaning counts as covered once one of its examples is
// verified (LEX-5, P2-09c).
func TestExampleCoverage(t *testing.T) {
	tx := dbtest.Tx(t)
	ctx := contributor(t, tx)
	svc := contributions.NewService(tx)
	sense := coreMeaning(ctx, t, tx, svc)

	added, err := svc.ProposeExample(ctx, sense, contributions.Example{Text: "Zz amenshi yatalala", TranslationEn: "The water is cold"})
	if err != nil || added.Trust.Status != "sourced" {
		t.Fatalf("propose example = %+v, %v", added, err)
	}
	if _, err := svc.ProposeExample(ctx, uuid.New(), contributions.Example{Text: "x", TranslationEn: "x"}); httpx.Status(err) != 404 {
		t.Errorf("unknown meaning: %v, want 404", err)
	}

	lex := lexicon.NewStore(tx)
	c, err := lex.ExampleCoverage(t.Context(), "tst")
	if err != nil || c.CoreSenses != 1 || c.WithExample != 0 || len(c.Missing) != 1 || c.Missing[0].Pending != 2 {
		t.Fatalf("coverage = %+v, %v; want 1 core meaning missing an example, 2 pending", c, err)
	}
	if _, err := tx.Exec(t.Context(), `UPDATE items SET status = 'verified' WHERE id = $1`, added.ExampleID); err != nil {
		t.Fatal(err)
	}
	if c, err := lex.ExampleCoverage(t.Context(), "tst"); err != nil || c.WithExample != 1 || len(c.Missing) != 0 {
		t.Errorf("after verifying an example: %+v, %v; want covered", c, err)
	}
}

// coreMeaning adds a verified meaning of a core concept in "tst", with one
// example waiting for review.
func coreMeaning(ctx context.Context, t *testing.T, tx pgx.Tx, svc *contributions.Service) uuid.UUID {
	t.Helper()
	if _, err := tx.Exec(t.Context(), `INSERT INTO languages (code, name, own_name, wave, status) VALUES ('tst', 'Test', 'Test', 9, 'active')`); err != nil {
		t.Fatal(err)
	}
	word := water
	word.Language = "tst"
	r, err := svc.Propose(ctx, word)
	if err != nil {
		t.Fatal(err)
	}
	// The meaning is verified and its concept on the core list; its example is not yet.
	for _, q := range []string{
		`UPDATE items SET status = 'verified' WHERE id = $1`,
		`UPDATE concepts SET core = true WHERE id = (SELECT concept_id FROM senses WHERE id = $1)`,
	} {
		if _, err := tx.Exec(t.Context(), q, r.SenseID); err != nil {
			t.Fatal(err)
		}
	}
	return r.SenseID
}
