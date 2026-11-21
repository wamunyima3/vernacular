package lexicon_test

import (
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

// TestCountEntries: of three entries, one is verified and one disputed, so
// two count in all and one as verified.
func TestCountEntries(t *testing.T) {
	tx := dbtest.Tx(t)
	ctx := t.Context()
	if _, err := tx.Exec(ctx, `INSERT INTO languages (code, name, own_name, wave, status) VALUES ('tst', 'Test', 'Test', 9, 'active')`); err != nil {
		t.Fatal(err)
	}
	claude, err := sources.Lookup(ctx, tx, "Claude")
	if err != nil {
		t.Fatal(err)
	}
	lex := lexicon.NewStore(tx)
	for i, w := range []struct{ headword, status string }{{"amenshi", "verified"}, {"umulilo", "disputed"}, {"ibwe", "machine"}} {
		added, err := lex.Add(ctx, lexicon.Proposal{
			Language: "tst", Headword: w.headword, GlossEn: "zz " + w.headword, Origin: items.Machine(claude, "m", string(rune('a'+i))),
		})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := tx.Exec(ctx, `UPDATE items SET status = $2 WHERE id = $1`, added.SenseID, w.status); err != nil {
			t.Fatal(err)
		}
	}
	counts, err := lex.CountEntries(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got := counts["tst"]; got != (lexicon.EntryCount{Verified: 1, Total: 2}) {
		t.Errorf("counts = %+v, want 1 verified of 2", got)
	}
}
