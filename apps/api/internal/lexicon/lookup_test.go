package lexicon_test

import (
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

// TestLookup: icipazzo (Bemba) and zzcinthu (Nyanja) share the concept
// "zz ichipazzo".
func TestLookup(t *testing.T) {
	lex, tx := searchFixture(t)
	ctx := t.Context()
	claude, err := sources.Lookup(ctx, tx, "Claude")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := lex.Add(ctx, lexicon.Proposal{Language: "nya", Headword: "zzcinthu", GlossEn: "zz ichipazzo", Origin: items.Machine(claude, "m", "n")}); err != nil {
		t.Fatal(err)
	}
	check := func(name string, hits []lexicon.LookupHit, err error, match, headword string) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
		for _, h := range hits {
			if h.GlossEn == "zz ichipazzo" {
				if h.Match != match || len(h.Entries) != 1 || h.Entries[0].Headword != headword {
					t.Errorf("%s: %s with %d entries (first %q), want %s with only %s", name, h.Match, len(h.Entries), h.Entries[0].Headword, match, headword)
				}
				return
			}
		}
		t.Errorf("%s: concept not found in %+v", name, hits)
	}
	hits, err := lex.FromEnglish(ctx, "ZZ Ichipazzo", []string{"nya"}, 5)
	check("from English", hits, err, "headword", "zzcinthu")
	hits, err = lex.FromEnglish(ctx, "zz ichipazo", []string{"bem"}, 5)
	check("typo in English", hits, err, "similar", "icipazzo")
	hits, err = lex.ToEnglish(ctx, "bem", "ichipazzo", 5)
	check("to English", hits, err, "variant", "icipazzo")
}
