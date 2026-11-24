package candidates

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"html"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// googleCodes maps our ISO 639-3 codes to Google Translate's codes. Only
// Bemba and Nyanja are offered among our languages.
var googleCodes = map[string]language.Tag{
	"bem": language.Make("bem"),
	"nya": language.Make("ny"),
}

// googleModel is Google's neural machine translation model.
const googleModel = "nmt"

// Google proposes the machine translation of each concept's gloss. It gives
// one word or phrase and no grammar, so part of speech stays empty for
// reviewers to fill in.
type Google struct {
	client *translate.Client
}

// NewGoogle wraps a Translation API (v2) client.
func NewGoogle(client *translate.Client) *Google { return &Google{client: client} }

// Generate translates one concept's English gloss.
func (g *Google) Generate(ctx context.Context, r Request) (Result, error) {
	model, text, err := g.translate(ctx, r.Language, r.Concept.GlossEn)
	if err != nil {
		return Result{}, err
	}
	return Result{Model: model, Candidates: []Candidate{{Headword: text}}}, nil
}

// Phrase translates one line of the phrase list.
func (g *Google) Phrase(ctx context.Context, r PhraseRequest) (PhraseResult, error) {
	model, text, err := g.translate(ctx, r.Language, r.SourceEn)
	if err != nil {
		return PhraseResult{}, err
	}
	return PhraseResult{Model: model, Candidates: []PhraseCandidate{{Text: text}}}, nil
}

// translate returns Google's model version and its translation of English text.
func (g *Google) translate(ctx context.Context, lang, text string) (model, out string, err error) {
	target, ok := googleCodes[lang]
	if !ok {
		return "", "", fmt.Errorf("google translate does not offer %q", lang)
	}
	res, err := g.client.Translate(ctx, []string{text}, target, &translate.Options{
		Source: language.English, Format: translate.Text, Model: googleModel,
	})
	if err != nil {
		return "", "", fmt.Errorf("google translate: %w", err)
	}
	if len(res) == 0 {
		return "", "", errors.New("google translate returned nothing")
	}
	return "google-translate-" + cmp.Or(res[0].Model, googleModel), html.UnescapeString(res[0].Text), nil
}
