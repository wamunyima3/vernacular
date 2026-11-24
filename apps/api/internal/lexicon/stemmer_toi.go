package lexicon

// Tonga (docs/style-guides/toi.md). Citation forms carry no augment, so one
// layer: the infinitive ku-/kw- or a noun class prefix: mu-/mw-, ba-, mi-,
// li-, ma-, ci-/cy-, zi-, lu-/lw-, ka-, tu-/tw-, bu-/bw-. Class 5 i- and the
// class 9/10 nasals are not stripped, because they are hard to tell from
// the stem itself; so ibbwe and mabwe do not share a stem yet.
//
// Words reach the stemmer after Standardise (ch → c) and folding.
func init() {
	stemmers["toi"] = newStemmer([]string{
		"kw", "ku",
		"mw", "mu",
		"ba",
		"mi",
		"li",
		"ma",
		"cy", "ci",
		"zi",
		"lw", "lu",
		"ka",
		"tw", "tu",
		"bw", "bu",
	})
}
