package lexicon_test

import (
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
)

func TestStandardise(t *testing.T) {
	tests := []struct{ lang, in, want, variant string }{
		{"bem", "ichibemba", "icibemba", "ichibemba"},
		{"bem", "icibemba", "icibemba", ""},
		{"bem", "ŋombe", "ng'ombe", "ŋombe"},
		{"nya", "chinthu", "cinthu", "chinthu"},
		{"nya", "aŵiri", "awiri", "aŵiri"},
		{"toi", "chintu", "cintu", "chintu"},
		{"toi", "ŋombe", "ng'ombe", "ŋombe"},
		{"loz", "kucha", "kuca", "kucha"},
		{"loz", "ng'aka", "ñaka", "ng'aka"},
		{"zzz", "chintu", "chintu", ""}, // no rules
	}
	for _, tt := range tests {
		got, variant := lexicon.Standardise(tt.lang, tt.in)
		if got != tt.want || variant != tt.variant {
			t.Errorf("Standardise(%s, %q) = (%q, %q), want (%q, %q)", tt.lang, tt.in, got, variant, tt.want, tt.variant)
		}
	}
}
