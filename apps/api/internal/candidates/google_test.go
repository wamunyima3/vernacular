package candidates

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/translate"
	"google.golang.org/api/option"
)

func fakeGoogle(t *testing.T, gotTarget *string) *Google {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		*gotTarget = r.Form.Get("target")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"translations":[{"translatedText":"menshi &amp; ifyakulya","model":"nmt"}]}}`))
	}))
	t.Cleanup(srv.Close)
	client, err := translate.NewClient(context.Background(), option.WithEndpoint(srv.URL), option.WithAPIKey("test"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close() })
	return NewGoogle(client)
}

func TestGoogleGenerate(t *testing.T) {
	var target string
	g := fakeGoogle(t, &target)

	got, err := g.Generate(context.Background(), Request{Language: "nya", Concept: water.Concept})
	if err != nil {
		t.Fatal(err)
	}
	if target != "ny" {
		t.Errorf("target = %q, want ny (Google's code for Nyanja)", target)
	}
	if got.Model != "google-translate-nmt" || len(got.Candidates) != 1 || got.Candidates[0].Headword != "menshi & ifyakulya" {
		t.Errorf("Generate = %+v", got)
	}

	if _, err := g.Generate(context.Background(), Request{Language: "loz", Concept: water.Concept}); err == nil {
		t.Error("Lozi is not offered by Google Translate; want an error")
	}
}

func TestGooglePhrase(t *testing.T) {
	var target string
	got, err := fakeGoogle(t, &target).Phrase(context.Background(), PhraseRequest{Language: "bem", SourceEn: "Good morning"})
	if err != nil {
		t.Fatal(err)
	}
	if target != "bem" || got.Model != "google-translate-nmt" || len(got.Candidates) != 1 || got.Candidates[0].Text != "menshi & ifyakulya" {
		t.Errorf("Phrase = %+v (target %q)", got, target)
	}
}
