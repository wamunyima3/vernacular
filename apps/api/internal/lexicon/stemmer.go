package lexicon

import (
	"slices"
	"strings"
)

// No stemmers exist for Zambian languages (architecture §10), so these are
// ours: light, rule-based prefix strippers for search (LEX-4). Bantu words
// inflect mostly at the front, so a stem is what is left after removing, in
// order, the layers a word can start with. Each layer removes at most one
// prefix, the longest that matches, and only if a stem of at least minStem
// letters remains. Stems are for matching, not linguistics: a match through a
// stem is reported as a `root` match and ranked below headword and variant
// matches.

// minStem is 2 because real roots can be that short (-ti tree, -pi palm).
// Search uses a stem match only when the stem has at least three letters
// (lexicon.MinSearchStem), since two-letter stems group unrelated words.
const minStem = 2

// MinSearchStem is the shortest stem that search matches on.
const MinSearchStem = 3

type stemmer struct {
	layers [][]string // prefix sets, applied in order; each sorted longest first
}

func newStemmer(layers ...[]string) stemmer {
	for _, l := range layers {
		slices.SortFunc(l, func(a, b string) int { return len(b) - len(a) })
	}
	return stemmer{layers: layers}
}

func (s stemmer) stem(word string) string {
	for _, layer := range s.layers {
		for _, p := range layer {
			if rest, ok := strings.CutPrefix(word, p); ok && len([]rune(rest)) >= minStem {
				word = rest
				break
			}
		}
	}
	return word
}

var stemmers = map[string]stemmer{}

// Stem returns the search stem of a folded word in a language. Languages
// without a stemmer return the word unchanged. Multi-word phrases are not
// stemmed.
func Stem(language, word string) string {
	s, ok := stemmers[language]
	if !ok || strings.Contains(word, " ") {
		return word
	}
	return s.stem(word)
}
