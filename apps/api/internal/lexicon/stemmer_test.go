package lexicon_test

import (
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
)

// Each group must share one stem, so a search for any form finds the others.
func TestStemGroupsShareAStem(t *testing.T) {
	groups := []struct {
		lang  string
		words []string
	}{
		{"bem", []string{"umuntu", "abantu", "muntu", "bantu", "icintu", "ifintu", "akantu", "utuntu", "ubuntu"}},
		{"bem", []string{"ukulya", "kulya"}},
		{"bem", []string{"umuti", "imiti", "muti"}},
		{"bem", []string{"ulupi", "lupi"}},
		{"nya", []string{"munthu", "anthu", "cinthu", "zinthu"}},
		{"nya", []string{"mtengo", "mitengo", "mutengo"}},
		{"nya", []string{"kudya", "dya"}},
		{"toi", []string{"muntu", "bantu", "cintu", "zintu", "kantu", "tuntu", "buntu"}},
		{"toi", []string{"musamu", "misamu"}},
		{"toi", []string{"kulya", "lya"}},
		{"loz", []string{"litapi", "matapi"}},
		{"loz", []string{"komu", "likomu"}},
	}
	for _, g := range groups {
		want := lexicon.Stem(g.lang, g.words[0])
		for _, w := range g.words[1:] {
			if got := lexicon.Stem(g.lang, w); got != want {
				t.Errorf("Stem(%s, %q) = %q, want %q (same as %q)", g.lang, w, got, want, g.words[0])
			}
		}
	}
}

func TestStemLeavesShortWordsPhrasesAndUnknownLanguages(t *testing.T) {
	tests := []struct{ lang, in, want string }{
		{"bem", "ama", "ama"},               // nothing would be left
		{"bem", "muli shani", "muli shani"}, // phrases are not stemmed
		{"bem", "ee", "ee"},
		{"zzz", "mutu", "mutu"},     // no stemmer for the language
		{"nya", "nyumba", "nyumba"}, // class 9 nasal stays
	}
	for _, tt := range tests {
		if got := lexicon.Stem(tt.lang, tt.in); got != tt.want {
			t.Errorf("Stem(%s, %q) = %q, want %q", tt.lang, tt.in, got, tt.want)
		}
	}
}
