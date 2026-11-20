package lexicon

// Nyanja (docs/style-guides/nya.md). Nyanja has no augment. One layer:
// the infinitive ku-/kw- or a noun class prefix: mu-/m-/mw-, a-, mi-, li-,
// ma-, ci-/cy-, zi-, ka-, ti-, u-. Class 9/10 nasals (n-, ny-) are not
// stripped, because they are hard to tell from the stem itself.
//
// Words reach the stemmer after Standardise (ch → c) and folding.
func init() {
	stemmers["nya"] = newStemmer([]string{
		"kw", "ku",
		"mw", "mu", "m",
		"a",
		"mi",
		"li",
		"ma",
		"cy", "ci",
		"zi",
		"ka",
		"ti",
		"u",
	})
}
