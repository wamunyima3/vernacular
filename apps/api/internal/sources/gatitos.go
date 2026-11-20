package sources

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"strings"
)

// GATITOS en_xx files are JSONL: {"src": "<English>", "trgs": ["<word>", ...]}.
// The parser only extracts; the lexicon normalises on write.
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
				if strings.TrimSpace(trg) == "" {
					continue
				}
				if !yield(Record{Text: trg, English: row.Src, Reference: ref}, nil) {
					return
				}
			}
		}
		if err := sc.Err(); err != nil {
			yield(Record{}, fmt.Errorf("%s: %w", g.file, err))
		}
	}
}
