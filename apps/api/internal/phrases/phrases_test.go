package phrases_test

import (
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/phrases"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
)

func TestPack(t *testing.T) {
	tx := dbtest.Tx(t)
	ctx := t.Context()
	if _, err := tx.Exec(ctx, `
		WITH i AS (INSERT INTO items (kind, language_code, status) VALUES ('phrase', 'bem', 'verified'), ('phrase', 'bem', 'machine')
			RETURNING id, status)
		INSERT INTO phrases (id, language_code, pack, text, source_en)
		SELECT id, 'bem', 'zz-greetings', CASE WHEN status = 'verified' THEN 'Mwashibukeni' ELSE 'Mulishani' END,
			CASE WHEN status = 'verified' THEN 'Good morning' ELSE 'How are you' END FROM i`); err != nil {
		t.Fatal(err)
	}
	store := phrases.NewStore(tx)
	for _, tt := range []struct {
		pack, min string
		want      []string
	}{
		{"zz-greetings", "verified", []string{"Mwashibukeni"}},
		{"zz-greetings", "machine", []string{"Mwashibukeni", "Mulishani"}},
		{"zz-numbers", "machine", nil},
	} {
		got, err := store.Pack(ctx, "bem", tt.pack, tt.min)
		if err != nil {
			t.Fatal(err)
		}
		var texts []string
		for _, p := range got {
			texts = append(texts, p.Text)
		}
		if len(texts) != len(tt.want) || (len(texts) > 0 && texts[0] != tt.want[0]) {
			t.Errorf("%s at %s = %v, want %v", tt.pack, tt.min, texts, tt.want)
		}
	}
}

// TestMatch: only the verified phrase matches, from English or from Bemba,
// ignoring case and accents.
func TestMatch(t *testing.T) {
	tx := dbtest.Tx(t)
	ctx := t.Context()
	if _, err := tx.Exec(ctx, `
		WITH i AS (INSERT INTO items (kind, language_code, status) VALUES ('phrase', 'bem', 'verified'), ('phrase', 'bem', 'sourced')
			RETURNING id, status)
		INSERT INTO phrases (id, language_code, pack, text, source_en)
		SELECT id, 'bem', 'zz-greetings', CASE WHEN status = 'verified' THEN 'Zz mwashibukeni' ELSE 'Zz mulishani' END,
			CASE WHEN status = 'verified' THEN 'Zz good morning' ELSE 'Zz how are you' END FROM i`); err != nil {
		t.Fatal(err)
	}
	store := phrases.NewStore(tx)
	for _, tt := range []struct {
		text        string
		fromEnglish bool
		want        int
	}{
		{"ZZ GOOD MÓRNING", true, 1},
		{"zz mwashibukeni", false, 1},
		{"Zz how are you", true, 0}, // only sourced
		{"zz good", true, 0},        // not exact
	} {
		got, err := store.Match(ctx, "bem", tt.text, tt.fromEnglish)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != tt.want {
			t.Errorf("Match(%q) = %d phrases, want %d", tt.text, len(got), tt.want)
		}
	}
}
