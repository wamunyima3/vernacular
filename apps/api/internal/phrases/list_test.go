package phrases_test

import (
	"strings"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/phrases"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
)

func TestImportList(t *testing.T) {
	tx := dbtest.Tx(t)
	const list = "pack,source_en,notes\nzz-greetings,Zz good morning,\nzz-greetings,Zz go well,said to someone leaving\n"
	for i, want := range []int{2, 0} { // re-running adds nothing
		added, err := phrases.ImportList(t.Context(), tx, strings.NewReader(list))
		if err != nil {
			t.Fatal(err)
		}
		if added != want {
			t.Errorf("run %d added %d, want %d", i+1, added, want)
		}
	}
	var gloss string
	var core bool
	if err := tx.QueryRow(t.Context(), `SELECT c.gloss_en, c.core FROM phrase_list l JOIN concepts c ON c.id = l.concept_id
		WHERE l.source_en = 'Zz go well'`).Scan(&gloss, &core); err != nil {
		t.Fatal(err)
	}
	if gloss != "zz go well" || core {
		t.Errorf("concept = %q (core %v), want zz go well, not core", gloss, core)
	}
	if _, err := phrases.ImportList(t.Context(), tx, strings.NewReader("pack,gloss\nx,y\n")); err == nil {
		t.Error("a wrong header was accepted")
	}
}
