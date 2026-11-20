package candidates

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/v2"

	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
)

// ClaudeModel is the model candidates are asked of. The model that actually
// answered (which differs when a refusal falls back) is stored per item.
const ClaudeModel = "claude-opus-5-5"

// answerSchema constrains the reply (structured outputs), so it always
// parses and every part of speech is one the database accepts.
var answerSchema = withPartsOfSpeech(mustSchema(`{
  "type": "object",
  "additionalProperties": false,
  "required": ["candidates"],
  "properties": {
    "candidates": {
      "type": "array",
      "items": {
        "type": "object",
        "additionalProperties": false,
        "required": ["headword", "part_of_speech", "noun_class", "plural"],
        "properties": {
          "headword": {"type": "string", "description": "Dictionary form, in the style guide's spelling"},
          "part_of_speech": {"type": "string"},
          "noun_class": {"type": "string", "description": "Singular/plural class pair such as 1/2 or 7/8; empty if not a noun"},
          "plural": {"type": "string", "description": "Plural form; empty if none"}
        }
      }
    }
  }
}`))

// withPartsOfSpeech fills the part_of_speech enum from the database's enum,
// so a new part of speech needs no change here.
func withPartsOfSpeech(schema map[string]any) map[string]any {
	item := schema["properties"].(map[string]any)["candidates"].(map[string]any)["items"].(map[string]any)
	pos := item["properties"].(map[string]any)["part_of_speech"].(map[string]any)
	values := lexicon.PartsOfSpeech()
	enum := make([]any, len(values))
	for i, v := range values {
		enum[i] = v
	}
	pos["enum"] = enum
	return schema
}

func mustSchema(s string) map[string]any {
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		panic(err)
	}
	return m
}

// Claude proposes candidates with the Claude API. Each language's system
// prompt holds its style guide and is cached, so only the concept is new
// input on each call.
type Claude struct {
	client  anthropic.Client
	systems map[string]string // language code → system prompt
}

// languageNames lists the languages Claude is asked about, by code. Own
// names follow the style guides (c, not ch), which the languages table
// predates.
var languageNames = map[string]string{"bem": "Bemba (Icibemba)", "nya": "Nyanja (Cinyanja)"}

// NewClaude builds a generator whose prompts include the style guides in dir
// (docs/style-guides: common.md plus <code>.md per language).
func NewClaude(client anthropic.Client, dir string) (*Claude, error) {
	common, err := os.ReadFile(filepath.Join(dir, "common.md"))
	if err != nil {
		return nil, fmt.Errorf("read style guide: %w", err)
	}
	systems := map[string]string{}
	for code, name := range languageNames {
		guide, err := os.ReadFile(filepath.Join(dir, code+".md"))
		if err != nil {
			return nil, fmt.Errorf("read style guide: %w", err)
		}
		systems[code] = systemPrompt(name, string(common), string(guide))
	}
	return &Claude{client: client, systems: systems}, nil
}

func systemPrompt(language, common, guide string) string {
	return `You propose dictionary entries in ` + language + `, a Zambian language, for an open dataset that public services, health workers and app developers rely on.

Each English concept you receive comes from a core vocabulary list. Give the ` + language + ` words or short expressions that speakers in Zambia use today for it, most common first, and at most three. Prefer the everyday Zambian form over a Malawian, archaic or Bible-only one.

What you return is stored labelled as machine-generated, and two independent native-speaker reviewers check it before anyone relies on it. A wrong word costs the reviewers time and can mislead them, so give only words you are confident about. Returning no candidates is the right answer when you are unsure.

Spell every word as the style guides below say.

<style_guide name="common">
` + common + `
</style_guide>

<style_guide name="` + language + `">
` + guide + `
</style_guide>`
}

func userPrompt(r Request) string {
	var b strings.Builder
	fmt.Fprintf(&b, "English concept: %q\nDomain: %s\n", r.Concept.GlossEn, r.Concept.Domain)
	if r.Concept.Definition != "" {
		fmt.Fprintf(&b, "Meaning: %s\n", r.Concept.Definition)
	}
	return b.String()
}

// Generate asks Claude for one concept in one language.
func (c *Claude) Generate(ctx context.Context, r Request) (Result, error) {
	system, ok := c.systems[r.Language]
	if !ok {
		return Result{}, fmt.Errorf("no style guide for language %q", r.Language)
	}
	msg, err := c.client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
		Model:     ClaudeModel,
		MaxTokens: 16000,
		System: []anthropic.BetaTextBlockParam{{
			Text:         system,
			CacheControl: anthropic.NewBetaCacheControlEphemeralParam(),
		}},
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(userPrompt(r))),
		},
		OutputConfig: anthropic.BetaOutputConfigParam{
			Effort: anthropic.BetaOutputConfigEffortHigh, // accuracy over cost: reviewers' time is scarcer
			Format: anthropic.BetaJSONOutputFormatParam{Schema: answerSchema},
		},
		// A policy decline is retried on a fallback model in the same call.
		Fallbacks: anthropic.BetaFallbacksParamOfDefault(),
		Betas:     []anthropic.AnthropicBeta{anthropic.AnthropicBetaServerSideFallback2026_07_01},
	})
	if err != nil {
		return Result{}, fmt.Errorf("claude: %w", err)
	}
	return parse(msg)
}

func parse(msg *anthropic.BetaMessage) (Result, error) {
	switch msg.StopReason {
	case anthropic.BetaStopReasonRefusal:
		return Result{}, ErrRefused
	case anthropic.BetaStopReasonEndTurn:
	default:
		return Result{}, fmt.Errorf("claude stopped early: %s", msg.StopReason)
	}
	var text string
	for _, block := range msg.Content {
		if block.Type == "text" {
			text += block.Text
		}
	}
	var answer struct {
		Candidates []Candidate `json:"candidates"`
	}
	if err := json.Unmarshal([]byte(text), &answer); err != nil {
		return Result{}, fmt.Errorf("parse claude answer: %w", err)
	}
	return Result{Model: msg.Model, Candidates: answer.Candidates}, nil
}
