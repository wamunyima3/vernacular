package phrases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/machine"
	"github.com/wamunyima3/vernacular/apps/api/internal/phrases"
)

// verified matches "Good morning" only.
type verified struct{}

func (verified) Pack(context.Context, string, string, string) ([]phrases.Phrase, error) {
	return nil, nil
}

func (verified) Match(_ context.Context, _, text string, _ bool) ([]phrases.Phrase, error) {
	if text == "Good morning" {
		return []phrases.Phrase{{Text: "Mwashibukeni", SourceEn: "Good morning"}}, nil
	}
	return nil, nil
}

type model struct {
	err   error
	calls int
}

func (m *model) Translate(context.Context, string, string, string) (machine.Result, error) {
	m.calls++
	return machine.Result{Translation: "Mulale bwino", Confidence: 0.4, Model: "m-1"}, m.err
}

func translate(t *testing.T, mt machine.Translator, text string, ask bool) gen.Translate200JSONResponse {
	t.Helper()
	h := phrases.NewHandler(verified{}, mt)
	res, err := h.Translate(t.Context(), gen.TranslateRequestObject{Body: &gen.TranslateJSONRequestBody{
		Text: text, From: "eng", To: "bem", Machine: &ask,
	}})
	if err != nil {
		t.Fatal(err)
	}
	return res.(gen.Translate200JSONResponse)
}

func TestMachineTranslation(t *testing.T) {
	m := &model{}
	if got := translate(t, m, "Good night", true).Machine; got == nil || got.Status != gen.MachineTranslationStatusMachine || got.Model != "m-1" {
		t.Fatalf("no verified match, machine asked: %+v", got)
	}
	if got := translate(t, m, "Good morning", true); len(got.Data) != 1 || got.Machine != nil {
		t.Fatalf("a verified match comes alone: %+v", got)
	}
	if got := translate(t, m, "Good night", false).Machine; got != nil {
		t.Fatalf("machine not asked for: %+v", got)
	}
	if m.calls != 1 {
		t.Fatalf("the model was called %d times, want only when needed", m.calls)
	}
	if got := translate(t, nil, "Good night", true).Machine; got != nil {
		t.Fatal("no service should mean no machine translation")
	}
	failing := &model{err: errors.New("down")}
	if got := translate(t, failing, "Good night", true); got.Machine != nil || got.Data == nil {
		t.Fatalf("a failing service should leave machine out, not fail: %+v", got)
	}
}
