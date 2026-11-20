package lexicon

// Bemba (docs/style-guides/bem.md). Layers, in order:
//
//  1. Infinitive or noun class prefix, with or without the augment:
//     uku-/ukw- (to …), umu-/mu-, aba-/ba-, imi-/mi-, ili-/li-, ama-/ma-,
//     ici-/ci-, ifi-/fi-, ulu-/lu-, aka-/ka-, utu-/tu-, ubu-/bu-, in-/im-.
//  2. Tense and aspect markers after a subject prefix are not stripped yet:
//     conjugated verbs are rare in a dictionary query.
//
// Folding has already lowered case and removed accents. Prefixes ending in a
// vowel also come in a pre-vowel form (ukw-, umw-, ici- → icy-), listed too.
func init() {
	stemmers["bem"] = newStemmer([]string{
		"ukw", "uku", "ku",
		"umw", "umu", "mw", "mu",
		"aba", "ab", "ba",
		"imi", "mi",
		"ili", "li",
		"ama", "am", "ma",
		"icy", "ici", "ci",
		"ify", "ifi", "fi",
		"ulw", "ulu", "lu",
		"aka", "ka",
		"utw", "utu", "tu",
		"ubw", "ubu", "bu",
		"im", "in",
	})
}
