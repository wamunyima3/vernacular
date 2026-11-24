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

// phraseSchema constrains a phrase reply.
var phraseSchema = mustSchema(`{
  "type": "object",
  "additionalProperties": false,
  "required": ["phrases"],
  "properties": {
    "phrases": {
      "type": "array",
      "items": {
        "type": "object",
        "additionalProperties": false,
        "required": ["text", "notes"],
        "properties": {
          "text": {"type": "string", "description": "The phrase, in the style guide's spelling"},
          "notes": {"type": "string", "description": "When to use this form rather than another, such as to one person or many; empty if there is only one"}
        }
      }
    }
  }
}`)

// Claude proposes candidates with the Claude API. Each language's system
// prompts hold its style guide and are cached, so only the concept or
// phrase is new input on each call.
type Claude struct {
	client  anthropic.Client
	words   map[string]string // language code → system prompt for words
	phrases map[string]string // language code → system prompt for phrases
}

// NewClaude builds a generator whose prompts include the style guides in dir
// (docs/style-guides: common.md plus <code>.md per language). Claude is asked
// about every language with a guide, by the name in the guide's title, so its
// own name follows the guide's spelling (c, not ch).
func NewClaude(client anthropic.Client, dir string) (*Claude, error) {
	common, err := os.ReadFile(filepath.Join(dir, "common.md"))
	if err != nil {
		return nil, fmt.Errorf("read style guide: %w", err)
	}
	paths, err := filepath.Glob(filepath.Join(dir, "[a-z][a-z][a-z].md"))
	if err != nil {
		return nil, err
	}
	c := &Claude{client: client, words: map[string]string{}, phrases: map[string]string{}}
	for _, path := range paths {
		guide, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read style guide: %w", err)
		}
		name, ok := guideName(guide)
		if !ok {
			return nil, fmt.Errorf("style guide %s: the first line must be \"# <Language> style guide\"", path)
		}
		code := strings.TrimSuffix(filepath.Base(path), ".md")
		c.words[code] = systemPrompt(name, wordTask(name), string(common), string(guide))
		c.phrases[code] = systemPrompt(name, phraseTask(name), string(common), string(guide))
	}
	return c, nil
}

// guideName reads the language from a guide's title, "# Bemba (Icibemba) style guide".
func guideName(guide []byte) (string, bool) {
	title, _, _ := strings.Cut(string(guide), "\n")
	name, ok := strings.CutSuffix(strings.TrimPrefix(title, "# "), " style guide")
	return name, ok && name != "" && strings.HasPrefix(title, "# ")
}

func wordTask(language string) string {
	return `You propose dictionary entries in ` + language + `, a Zambian language, for an open dataset that public services, health workers and app developers rely on.

Each English concept you receive comes from a core vocabulary list. Give the ` + language + ` words or short expressions that speakers in Zambia use today for it, most common first, and at most three. Prefer the everyday Zambian form over a Malawian, archaic or Bible-only one.`
}

func phraseTask(language string) string {
	return `You translate short everyday phrases into ` + language + `, a Zambian language, for phrase packs that apps and public services show to the people they serve: greetings, numbers, dates, interface labels, health and money.

Each English line you receive names its pack and may carry a note on its context, such as "button". Give the ` + language + ` phrase speakers in Zambia would expect to see there, as natural speech rather than a word-for-word rendering. Give a second form only when the choice depends on who is addressed, such as one person or many, or an elder, and say so in its notes; at most two. Keep any placeholder in braces, such as {name}, exactly as written.`
}

func systemPrompt(language, task, common, guide string) string {
	return task + `

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

func phrasePrompt(r PhraseRequest) string {
	prompt := fmt.Sprintf("Pack: %s\nEnglish: %q\n", r.Pack, r.SourceEn)
	if r.Notes != "" {
		prompt += "Context: " + r.Notes + "\n"
	}
	return prompt
}

// Generate asks Claude for one concept in one language.
func (c *Claude) Generate(ctx context.Context, r Request) (Result, error) {
	var answer struct {
		Candidates []Candidate `json:"candidates"`
	}
	model, err := c.ask(ctx, c.words, r.Language, userPrompt(r), answerSchema, &answer)
	return Result{Model: model, Candidates: answer.Candidates}, err
}

// Phrase asks Claude for one line of the phrase list in one language.
func (c *Claude) Phrase(ctx context.Context, r PhraseRequest) (PhraseResult, error) {
	var answer struct {
		Phrases []PhraseCandidate `json:"phrases"`
	}
	model, err := c.ask(ctx, c.phrases, r.Language, phrasePrompt(r), phraseSchema, &answer)
	return PhraseResult{Model: model, Candidates: answer.Phrases}, err
}

// ask sends one prompt under the language's cached system prompt and decodes
// the structured answer into out. It returns the model that answered.
func (c *Claude) ask(ctx context.Context, systems map[string]string, language, prompt string, schema map[string]any, out any) (string, error) {
	system, ok := systems[language]
	if !ok {
		return "", fmt.Errorf("no style guide for language %q", language)
	}
	msg, err := c.client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
		Model:     ClaudeModel,
		MaxTokens: 16000,
		System: []anthropic.BetaTextBlockParam{{
			Text:         system,
			CacheControl: anthropic.NewBetaCacheControlEphemeralParam(),
		}},
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(prompt)),
		},
		OutputConfig: anthropic.BetaOutputConfigParam{
			Effort: anthropic.BetaOutputConfigEffortHigh, // accuracy over cost: reviewers' time is scarcer
			Format: anthropic.BetaJSONOutputFormatParam{Schema: schema},
		},
		// A policy decline is retried on a fallback model in the same call.
		Fallbacks: anthropic.BetaFallbacksParamOfDefault(),
		Betas:     []anthropic.AnthropicBeta{anthropic.AnthropicBetaServerSideFallback2026_07_01},
	})
	if err != nil {
		return "", fmt.Errorf("claude: %w", err)
	}
	return parse(msg, out)
}

// parse decodes a structured answer into out and returns the model that gave it.
func parse(msg *anthropic.BetaMessage, out any) (string, error) {
	switch msg.StopReason {
	case anthropic.BetaStopReasonRefusal:
		return "", ErrRefused
	case anthropic.BetaStopReasonEndTurn:
	default:
		return "", fmt.Errorf("claude stopped early: %s", msg.StopReason)
	}
	var text string
	for _, block := range msg.Content {
		if block.Type == "text" {
			text += block.Text
		}
	}
	if err := json.Unmarshal([]byte(text), out); err != nil {
		return "", fmt.Errorf("parse claude answer: %w", err)
	}
	return msg.Model, nil
}
