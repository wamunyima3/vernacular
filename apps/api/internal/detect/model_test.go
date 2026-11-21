package detect

import "testing"

var samples = map[string][]string{
	"bem": {"umuntu", "abantu", "amenshi", "icintu", "ukulya", "umulilo", "mwashibukeni mukwai", "natotela sana", "nshilefwaya", "ifyakulya fyaliwama"},
	"nya": {"munthu", "anthu", "madzi", "cinthu", "kudya", "moto", "mwadzuka bwanji", "zikomo kwambiri", "sindikufuna", "chakudya ndi chabwino"},
	"eng": {"person", "people", "water", "thing", "to eat", "fire", "good morning", "thank you very much", "I do not want", "the food is good"},
	"xxx": {"too", "little"},
}

func TestDetect(t *testing.T) {
	m := Train(samples, 30)
	if got := m.Languages(); len(got) != 3 {
		t.Fatalf("languages = %v, want bem, eng, nya (xxx has too little text)", got)
	}
	for text, want := range map[string]string{
		"Natotela sana mukwai":  "bem",
		"Zikomo kwambiri anthu": "nya",
		"Thank you, people":     "eng",
	} {
		got := m.Detect(text)
		if len(got) != 3 || got[0].Language != want {
			t.Errorf("Detect(%q) = %+v, want %s first", text, got, want)
			continue
		}
		sum := 0.0
		for _, s := range got {
			sum += s.Confidence
		}
		if sum < 0.999 || sum > 1.001 || got[0].Confidence <= got[1].Confidence {
			t.Errorf("Detect(%q) confidences %+v do not rank and sum to 1", text, got)
		}
	}
	if got := m.Detect("123 !!"); got != nil {
		t.Errorf("no letters: %+v, want nothing", got)
	}
}
