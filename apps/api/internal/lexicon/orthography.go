package lexicon

import "strings"

// orthography maps each language to the replacements that turn common
// spellings into the standard one from its style guide (docs/style-guides).
// The original spelling is kept as a variant, so search still finds it.
var orthography = map[string]*strings.Replacer{
	// bem.md: c for /tʃ/, ng' for /ŋ/.
	"bem": strings.NewReplacer("ch", "c", "ŋ", "ng'"),
	// nya.md: c for /tʃ/, w not Malawian ŵ, ng' for /ŋ/.
	"nya": strings.NewReplacer("ch", "c", "ŵ", "w", "ŋ", "ng'"),
	// toi.md: c for /tʃ/, ng' for /ŋ/. The bb/b rule is left to reviewers.
	"toi": strings.NewReplacer("ch", "c", "ŋ", "ng'"),
	// loz.md: c for /tʃ/, ñ for /ŋ/.
	"loz": strings.NewReplacer("ch", "c", "ng'", "ñ", "ŋ", "ñ"),
}

// Standardise returns the standard spelling of a normalised headword, and the
// input as a variant when it differs.
func Standardise(language, headword string) (standard string, variant string) {
	r, ok := orthography[language]
	if !ok {
		return headword, ""
	}
	standard = r.Replace(headword)
	if standard == headword {
		return headword, ""
	}
	return standard, headword
}
