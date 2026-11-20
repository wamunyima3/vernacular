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
	target, ok := googleCodes[r.Language]
	if !ok {
		return Result{}, fmt.Errorf("google translate does not offer %q", r.Language)
	}
	out, err := g.client.Translate(ctx, []string{r.Concept.GlossEn}, target, &translate.Options{
		Source: language.English, Format: translate.Text, Model: googleModel,
	})
	if err != nil {
		return Result{}, fmt.Errorf("google translate: %w", err)
	}
	if len(out) == 0 {
		return Result{}, errors.New("google translate returned nothing")
	}
	model := "google-translate-" + cmp.Or(out[0].Model, googleModel)
	return Result{Model: model, Candidates: []Candidate{{Headword: html.UnescapeString(out[0].Text)}}}, nil
}
