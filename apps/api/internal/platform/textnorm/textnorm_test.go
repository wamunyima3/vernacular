package textnorm_test

import (
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/textnorm"
)

// Shared by both tests so Go and SQL folding are checked on the same inputs.
var foldCases = map[string]string{
	"Amenshi":   "amenshi",
	"Ichíbemba": "ichibemba",
	"ichí":     "ichi",
	"ŋombe":     "ng'ombe",
	"Ŋombe":     "ng'ombe",
	"NG’OMBE":   "ng'ombe",
	"ngʼombe":   "ng'ombe",
	"aŵiri":     "awiri",
	"Mwâpoleni": "mwapoleni",
	"cinthu":    "cinthu",
	"":          "",
}

func TestFold(t *testing.T) {
	for in, want := range foldCases {
		if got := textnorm.Fold(in); got != want {
			t.Errorf("Fold(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNFC(t *testing.T) {
	if got := textnorm.NFC("icíntu"); got != "icíntu" {
		t.Errorf("NFC = %q, want precomposed icíntu", got)
	}
}

func TestFoldMatchesSQL(t *testing.T) {
	tx := dbtest.Tx(t)
	for in := range foldCases {
		var sql string
		if err := tx.QueryRow(t.Context(), "SELECT fold($1)", in).Scan(&sql); err != nil {
			t.Fatal(err)
		}
		if got := textnorm.Fold(in); got != sql {
			t.Errorf("Fold(%q): Go %q, SQL %q", in, got, sql)
		}
	}
}
