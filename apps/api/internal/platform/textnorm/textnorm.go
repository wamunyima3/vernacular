// Package textnorm normalises text on input and folds it for search (NFR-7).
//
// Fold must give the same result as the SQL function fold() in
// db/migrations/00003_core.sql; TestFoldMatchesSQL checks that.
package textnorm

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// apostrophes maps the curly and modifier apostrophes people type to the
// straight one the style guides use, and ŋ to its standard spelling ng'.
var apostrophes = strings.NewReplacer("’", "'", "ʼ", "'", "‘", "'", "ŋ", "ng'")

// NFC returns s in Unicode normal form C. Every stored string goes
// through it, so equal text is equal bytes.
func NFC(s string) string { return norm.NFC.String(s) }

// Fold returns the search key for s: NFC, lower case, straight apostrophes,
// ŋ as ng', and accents removed.
func Fold(s string) string {
	s = apostrophes.Replace(strings.ToLower(NFC(s)))
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	out, _, err := transform.String(t, s)
	if err != nil {
		return s // transform only fails on invalid input; searching the raw text is the safe fallback
	}
	return out
}
