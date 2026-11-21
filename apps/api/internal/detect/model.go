// Package detect identifies which supported language a text is in (TRN-1,
// ADR-0010), with a small character n-gram model trained on our own data.
package detect

import (
	"cmp"
	"maps"
	"math"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/wamunyima3/vernacular/apps/api/internal/platform/textnorm"
)

// maxN is the longest character n-gram the model counts.
const maxN = 4

// temperature scales the per-n-gram mean log-likelihood gap before the
// softmax, so confidences reflect how distinct a text is rather than how long.
const temperature = 4.0

// Score is how likely a text is in a language, from 0 to 1.
type Score struct {
	Language   string
	Confidence float64
}

// profile holds a language's smoothed log-probabilities, computed once.
type profile struct {
	logp   map[string]float64
	unseen float64 // log-probability of an n-gram the language never had
}

// Model is a multinomial naive Bayes classifier over character n-grams.
// It is read-only after Train, so safe for concurrent use.
type Model struct {
	langs map[string]profile
}

// Train builds a model from each language's sample texts. Languages with
// fewer than minChars characters of text are left out rather than guessed.
func Train(samples map[string][]string, minChars int) *Model {
	counts := map[string]map[string]int{}
	totals := map[string]int{}
	vocab := map[string]struct{}{}
	for lang, texts := range samples {
		if runeCount(texts) < minChars {
			continue
		}
		c := map[string]int{}
		for _, t := range texts {
			for _, g := range ngrams(t) {
				c[g]++
				totals[lang]++
				vocab[g] = struct{}{}
			}
		}
		counts[lang] = c
	}
	m := &Model{langs: make(map[string]profile, len(counts))}
	for lang, c := range counts {
		denom := math.Log(float64(totals[lang] + len(vocab) + 1))
		p := profile{logp: make(map[string]float64, len(c)), unseen: -denom}
		for g, n := range c {
			p.logp[g] = math.Log(float64(n+1)) - denom
		}
		m.langs[lang] = p
	}
	return m
}

// Languages lists the languages the model can detect, sorted.
func (m *Model) Languages() []string {
	return slices.Sorted(maps.Keys(m.langs))
}

// Detect ranks the model's languages for a text, most likely first. It is
// empty when the text has no letters or the model knows no language.
func (m *Model) Detect(text string) []Score {
	grams := ngrams(text)
	if len(grams) == 0 || len(m.langs) == 0 {
		return nil
	}
	langs := m.Languages()
	logs := make([]float64, len(langs))
	best := math.Inf(-1)
	for i, lang := range langs {
		p := m.langs[lang]
		for _, g := range grams {
			if lp, ok := p.logp[g]; ok {
				logs[i] += lp
			} else {
				logs[i] += p.unseen
			}
		}
		best = max(best, logs[i])
	}
	scores := make([]Score, len(langs))
	sum := 0.0
	for i, lang := range langs {
		c := math.Exp((logs[i] - best) / float64(len(grams)) * temperature)
		scores[i] = Score{Language: lang, Confidence: c}
		sum += c
	}
	for i := range scores {
		scores[i].Confidence /= sum
	}
	slices.SortStableFunc(scores, func(a, b Score) int { return cmp.Compare(b.Confidence, a.Confidence) })
	return scores
}

// ngrams folds the text (NFR-7), keeps letters and apostrophes as words, and
// returns every 1..maxN-character n-gram of each word padded with spaces.
func ngrams(text string) []string {
	words := strings.FieldsFunc(textnorm.Fold(text), func(r rune) bool { return !unicode.IsLetter(r) && r != '\'' })
	var out []string
	for _, w := range words {
		rs := []rune(" " + w + " ")
		for n := 1; n <= maxN; n++ {
			for i := 0; i+n <= len(rs); i++ {
				if g := string(rs[i : i+n]); g != " " {
					out = append(out, g)
				}
			}
		}
	}
	return out
}

func runeCount(texts []string) int {
	n := 0
	for _, t := range texts {
		n += utf8.RuneCountInString(t)
	}
	return n
}
