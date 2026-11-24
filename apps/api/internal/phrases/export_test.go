package phrases_test

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/phrases"
)

var pack = []phrases.Phrase{
	{ID: uuid.New(), SourceEn: "Good morning", Text: "Mwashibukeni", Notes: "To one or many\nPolite", Trust: items.Trust{Status: "verified"}},
	{ID: uuid.New(), SourceEn: "Good morning", Text: "Shibukeni", Trust: items.Trust{Status: "verified"}},
	{ID: uuid.New(), SourceEn: `Say "yes"`, Text: "Tila \"ee\"\\", Trust: items.Trust{Status: "reviewed"}},
}

// TestExportFormats: each format keys by English, keeps the first
// translation of a line, and escapes what its syntax needs (PHR-3, US-12).
func TestExportFormats(t *testing.T) {
	if m := phrases.Map(pack); len(m) != 2 || m["Good morning"] != "Mwashibukeni" {
		t.Errorf("map = %v, want 2 lines with the first translation", m)
	}

	po := string(phrases.PO("greetings", "bem", pack))
	for _, want := range []string{
		`"Language: bem\n"`,
		"# To one or many\n# Polite\n#. id: " + pack[0].ID.String() + ", status: verified\nmsgid \"Good morning\"\nmsgstr \"Mwashibukeni\"\n",
		`msgid "Say \"yes\""` + "\n" + `msgstr "Tila \"ee\"\\"`,
	} {
		if !strings.Contains(po, want) {
			t.Errorf("po lacks %q:\n%s", want, po)
		}
	}
	if strings.Count(po, `msgid "Good morning"`) != 1 {
		t.Errorf("po repeats a msgid:\n%s", po)
	}

	out, err := phrases.XLIFF("greetings", "bem", pack)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		TrgLang string `xml:"trgLang,attr"`
		Units   []struct {
			ID      string `xml:"id,attr"`
			Segment struct {
				State  string `xml:"state,attr"`
				Target string `xml:"target"`
			} `xml:"segment"`
		} `xml:"file>unit"`
	}
	if err := xml.Unmarshal(out, &doc); err != nil {
		t.Fatalf("xliff does not parse: %v\n%s", err, out)
	}
	if doc.TrgLang != "bem" || len(doc.Units) != 2 || doc.Units[0].Segment.State != "final" ||
		doc.Units[1].Segment.State != "reviewed" || doc.Units[1].Segment.Target != pack[2].Text {
		t.Errorf("xliff = %+v", doc)
	}
}
