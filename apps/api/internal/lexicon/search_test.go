package lexicon_test

import (
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

// TestSearch uses made-up Bemba words, so the stemmer and style-guide
// spelling apply without colliding with real entries.
func TestSearch(t *testing.T) {
	lex := searchFixture(t)
	ctx := t.Context()
	tests := []struct{ query, headword, match string }{
		{"UMUKÁZZO", "umukazzo", "headword"},
		{"ichipazzo", "icipazzo", "variant"},
		{"abakazzo", "umukazzo", "root"},
		{"umukazo", "umukazzo", "similar"},
		{"umupazzo", "umupazzo", ""}, // only a disputed sense: never returned
	}
	for _, tt := range tests {
		hits, err := lex.Search(ctx, "bem", tt.query, 50)
		if err != nil {
			t.Fatal(err)
		}
		h := find(hits, tt.headword)
		switch {
		case tt.match == "" && h != nil:
			t.Errorf("%q: %s returned, want left out", tt.query, tt.headword)
		case tt.match == "":
		case h == nil:
			t.Errorf("%q: %s not found", tt.query, tt.headword)
		case h.Match != tt.match || len(h.Entry.Senses) != 1 || h.Entry.Senses[0].Trust.Status != "machine" || len(h.Entry.Trust.Sources) != 1:
			t.Errorf("%q: hit = %s with %+v, want %s with one machine sense from Claude", tt.query, h.Match, h.Entry.Senses, tt.match)
		}
	}
}

func find(hits []lexicon.Hit, headword string) *lexicon.Hit {
	for i := range hits {
		if hits[i].Entry.Headword == headword {
			return &hits[i]
		}
	}
	return nil
}

func searchFixture(t *testing.T) *lexicon.Store {
	t.Helper()
	tx := dbtest.Tx(t)
	ctx := t.Context()
	claude, err := sources.Lookup(ctx, tx, "Claude")
	if err != nil {
		t.Fatal(err)
	}
	lex := lexicon.NewStore(tx)
	for i, w := range []struct{ headword, status string }{
		{"umukazzo", "machine"},
		{"ichipazzo", "machine"}, // stored as icipazzo, with ichipazzo as a variant
		{"umupazzo", "disputed"},
	} {
		added, err := lex.Add(ctx, lexicon.Proposal{
			Language: "bem", Headword: w.headword, GlossEn: "zz " + w.headword, Origin: items.Machine(claude, "m", string(rune('a'+i))),
		})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := tx.Exec(ctx, `UPDATE items SET status = $2 WHERE id = $1`, added.SenseID, w.status); err != nil {
			t.Fatal(err)
		}
	}

	return lex
}
