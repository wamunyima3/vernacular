package lexicon

// Lozi (docs/style-guides/loz.md). Lozi has no augment, so one layer: the
// infinitive ku-/kw- or a noun class prefix: mu-/mw-, ba-, mi-, li-, ma-,
// si-, lu-/lw-, bu-/bw-. Class 9 nouns have no prefix in the singular, and
// their li- plurals strip to the singular.
//
// Words reach the stemmer after Standardise (ch → c) and folding.
func init() {
	stemmers["loz"] = newStemmer([]string{
		"kw", "ku",
		"mw", "mu",
		"ba",
		"mi",
		"li",
		"ma",
		"si",
		"lw", "lu",
		"bw", "bu",
	})
}
