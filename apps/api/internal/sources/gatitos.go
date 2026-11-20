package sources

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"strings"
	"unicode"
	"unicode/utf8"
)

// GATITOS en_xx files are JSONL: {"src": "<English>", "trgs": ["<word>", ...]}.
type gatitos struct {
	file string
}

func (g gatitos) Parse(r io.Reader) iter.Seq2[Record, error] {
	return func(yield func(Record, error) bool) {
		sc := bufio.NewScanner(r)
		for line := 1; sc.Scan(); line++ {
			var row struct {
				Src  string   `json:"src"`
				Trgs []string `json:"trgs"`
			}
			if err := json.Unmarshal(sc.Bytes(), &row); err != nil {
				yield(Record{}, fmt.Errorf("%s:%d: %w", g.file, line, err))
				return
			}
			ref := fmt.Sprintf("%s:%d", g.file, line)
			for _, trg := range row.Trgs {
				text, english := headword(trg), gloss(row.Src)
				if text == "" || english == "" {
					continue
				}
				if !yield(Record{Text: text, English: english, Reference: ref}, nil) {
					return
				}
			}
		}
		if err := sc.Err(); err != nil {
			yield(Record{}, fmt.Errorf("%s: %w", g.file, err))
		}
	}
}

// headword applies the style guide: trimmed, single spaces, lower case.
// Proper nouns are lowered too; reviewers restore them.
func headword(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

// gloss lowers a sentence-initial capital ("Hello" → "hello") but keeps "I"
// and acronyms ("HIV"), so the same English meaning maps to one concept.
func gloss(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	first, size := utf8.DecodeRuneInString(s)
	word, _, _ := strings.Cut(s, " ")
	if size == 0 || word == "I" || strings.ToUpper(word) == word {
		return s
	}
	return string(unicode.ToLower(first)) + s[size:]
}
