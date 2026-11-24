package machine_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/machine"
)

func service(t *testing.T, status int, answer string) machine.Translator {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var in map[string]string
		if r.URL.Path != "/translate" || json.NewDecoder(r.Body).Decode(&in) != nil || in["text"] != "Good night" || in["to"] != "bem" {
			t.Errorf("request %s %v", r.URL.Path, in)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(answer))
	}))
	t.Cleanup(srv.Close)
	return machine.New(srv.URL + "/")
}

func TestClient(t *testing.T) {
	if machine.New("") != nil {
		t.Fatal("no URL should mean no translator")
	}
	r, err := service(t, 200, `{"translation":"Mulale bwino","confidence":0.42,"model":"m-1"}`).
		Translate(t.Context(), "Good night", "eng", "bem")
	if err != nil || r != (machine.Result{Translation: "Mulale bwino", Confidence: 0.42, Model: "m-1"}) {
		t.Fatalf("Translate = %+v, %v", r, err)
	}
	failures := map[string]struct {
		status int
		answer string
	}{
		"unsupported": {422, `{}`},
		"down":        {503, ``},
		"no model":    {200, `{"translation":"x","confidence":0.5}`},
		"confidence":  {200, `{"translation":"x","confidence":2,"model":"m"}`},
	}
	for name, f := range failures {
		_, err := service(t, f.status, f.answer).Translate(t.Context(), "Good night", "eng", "bem")
		if err == nil || (name == "unsupported") != errors.Is(err, machine.ErrUnsupported) {
			t.Errorf("%s: %v", name, err)
		}
	}
}
