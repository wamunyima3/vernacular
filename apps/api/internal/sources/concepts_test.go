package sources_test

import (
	"strings"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

func TestImportConcepts(t *testing.T) {
	tx := dbtest.Tx(t)
	csv := "domain,gloss_en\nhealth,Zz clinic\nmoney,zz kwacha\nhealth,zz clinic\n"

	added, err := sources.ImportConcepts(t.Context(), tx, strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}
	if added != 2 {
		t.Errorf("added = %d, want 2 (the duplicate counts once)", added)
	}
	var domain string
	var core bool
	if err := tx.QueryRow(t.Context(), "SELECT domain, core FROM concepts WHERE gloss_en = 'zz clinic'").Scan(&domain, &core); err != nil {
		t.Fatal(err)
	}
	if domain != "health" || !core {
		t.Errorf("zz clinic = (%s, %v), want (health, true)", domain, core)
	}
}

func TestImportConceptsRejectsBadInput(t *testing.T) {
	for name, csv := range map[string]string{
		"wrong header":  "gloss,domain\nclinic,health\n",
		"missing gloss": "domain,gloss_en\nhealth,\n",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := sources.ImportConcepts(t.Context(), dbtest.Tx(t), strings.NewReader(csv)); err == nil {
				t.Error("import succeeded, want an error")
			}
		})
	}
}
