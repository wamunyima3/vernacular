package candidates

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"

	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
)

// fakeClaude answers like the Messages API and keeps the last request.
func fakeClaude(t *testing.T, stopReason, text string) (*Claude, *map[string]any, *http.Header) {
	t.Helper()
	var body map[string]any
	var header http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Errorf("request body: %v", err)
		}
		header = r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "msg_test", "type": "message", "role": "assistant", "model": "claude-opus-5-5",
			"stop_reason": stopReason, "content": []any{map[string]any{"type": "text", "text": text}},
			"usage": map[string]any{"input_tokens": 1, "output_tokens": 1},
		})
	}))
	t.Cleanup(srv.Close)

	dir := t.TempDir()
	for f, title := range map[string]string{"common.md": "Common", "bem.md": "Bemba (Icibemba)", "nya.md": "Nyanja (Cinyanja)"} {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("# "+title+" style guide\nUse c, not ch."), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	client := anthropic.NewClient(option.WithBaseURL(srv.URL), option.WithAPIKey("test"), option.WithMaxRetries(0))
	c, err := NewClaude(client, dir)
	if err != nil {
		t.Fatal(err)
	}
	return c, &body, &header
}

var water = Request{Language: "bem", Concept: lexicon.Concept{GlossEn: "water", Domain: "everyday"}}

func TestClaudeGenerate(t *testing.T) {
	c, body, header := fakeClaude(t, "end_turn",
		`{"candidates":[{"headword":"amenshi","part_of_speech":"noun","noun_class":"6","plural":""}]}`)

	got, err := c.Generate(context.Background(), water)
	if err != nil {
		t.Fatal(err)
	}
	want := Result{Model: "claude-opus-5-5", Candidates: []Candidate{{Headword: "amenshi", PartOfSpeech: "noun", NounClass: "6"}}}
	if got.Model != want.Model || len(got.Candidates) != 1 || got.Candidates[0] != want.Candidates[0] {
		t.Errorf("Generate = %+v, want %+v", got, want)
	}

	req := *body
	if req["model"] != ClaudeModel {
		t.Errorf("model = %v", req["model"])
	}
	system := req["system"].([]any)[0].(map[string]any)
	if system["cache_control"] == nil {
		t.Error("the system prompt (with the style guide) must be cached")
	}
	out := req["output_config"].(map[string]any)
	if out["format"].(map[string]any)["type"] != "json_schema" || out["effort"] != "high" {
		t.Errorf("output_config = %v", out)
	}
	pos := out["format"].(map[string]any)["schema"].(map[string]any)["properties"].(map[string]any)["candidates"].(map[string]any)["items"].(map[string]any)["properties"].(map[string]any)["part_of_speech"].(map[string]any)
	if enum, _ := pos["enum"].([]any); len(enum) != len(lexicon.PartsOfSpeech()) {
		t.Errorf("part_of_speech enum = %v, want the database's %d values", pos["enum"], len(lexicon.PartsOfSpeech()))
	}
	if req["fallbacks"] != "default" || header.Get("anthropic-beta") != anthropic.AnthropicBetaServerSideFallback2026_07_01 {
		t.Errorf("fallbacks = %v, beta header = %q", req["fallbacks"], header.Get("anthropic-beta"))
	}
}

func TestClaudeGenerateFailures(t *testing.T) {
	tests := []struct {
		name, stop, text string
		req              Request
		wantErr          error
	}{
		{"refusal is final", "refusal", "", water, ErrRefused},
		{"cut off", "max_tokens", `{"candidates":[`, water, nil},
		{"no style guide", "end_turn", "{}", Request{Language: "loz"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _, _ := fakeClaude(t, tt.stop, tt.text)
			_, err := c.Generate(context.Background(), tt.req)
			if err == nil {
				t.Fatal("want an error")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestClaudePhrase(t *testing.T) {
	c, body, _ := fakeClaude(t, "end_turn", `{"phrases":[{"text":"Mwashibukeni","notes":"to many or an elder"}]}`)
	got, err := c.Phrase(context.Background(), PhraseRequest{Language: "bem", Pack: "greetings", SourceEn: "Good morning", Notes: "said on arrival"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Model != "claude-opus-5-5" || len(got.Candidates) != 1 || got.Candidates[0].Text != "Mwashibukeni" {
		t.Errorf("Phrase = %+v", got)
	}
	req := *body
	system := req["system"].([]any)[0].(map[string]any)["text"].(string)
	user := req["messages"].([]any)[0].(map[string]any)["content"].([]any)[0].(map[string]any)["text"].(string)
	if !strings.Contains(system, "phrase packs") || !strings.Contains(system, "{name}") ||
		!strings.Contains(user, `"Good morning"`) || !strings.Contains(user, "said on arrival") {
		t.Errorf("system %q\nuser %q", system, user)
	}
}

// TestClaudeReadsTheStyleGuides: every guide in docs/style-guides has a title
// Claude can name its language by.
func TestClaudeReadsTheStyleGuides(t *testing.T) {
	c, err := NewClaude(anthropic.NewClient(option.WithAPIKey("test")), "../../../../docs/style-guides")
	if err != nil {
		t.Fatal(err)
	}
	for _, code := range []string{"bem", "nya", "toi"} {
		if !strings.Contains(c.words[code], "style guide") {
			t.Errorf("no prompt for %s", code)
		}
	}
}
