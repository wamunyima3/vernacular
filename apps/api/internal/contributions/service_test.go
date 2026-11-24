package contributions_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/contributions"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

var water = contributions.Input{
	Language: "bem", Headword: "zzamenshi", PartOfSpeech: "noun", NounClass: "6", GlossEn: "zz water",
	Example: &contributions.Example{Text: "Ndefwaya  zzamenshi.", TranslationEn: "I want water."},
}

// contributor adds user c1 and signs them in.
func contributor(t *testing.T, tx pgx.Tx) context.Context {
	t.Helper()
	if _, err := tx.Exec(t.Context(), `INSERT INTO users (id, name, email, email_verified, role)
		VALUES ('c1', 'Chanda', 'chanda@test.invalid', true, 'contributor')`); err != nil {
		t.Fatal(err)
	}
	return accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "c1", Role: gen.RoleContributor})
}

func TestPropose(t *testing.T) {
	tx := dbtest.Tx(t)
	ctx := contributor(t, tx)
	svc := contributions.NewService(tx)

	r, err := svc.Propose(ctx, water)
	if err != nil {
		t.Fatal(err)
	}
	if !r.NewSense || r.ExampleID == nil {
		t.Errorf("receipt = %+v, want a new sense and an example", r)
	}
	if r.Trust.Status != "sourced" || len(r.Trust.Sources) != 1 || r.Trust.Sources[0].Name != contributions.SourceName {
		t.Errorf("trust = %+v, want sourced, citing %s", r.Trust, contributions.SourceName)
	}
	var author, text string
	if err := tx.QueryRow(t.Context(), `SELECT i.author_id, x.text FROM examples x JOIN items i ON i.id = x.id WHERE x.id = $1`,
		*r.ExampleID).Scan(&author, &text); err != nil {
		t.Fatal(err)
	}
	if author != "c1" || text != "Ndefwaya zzamenshi." {
		t.Errorf("example author %q text %q, want c1 and squashed spaces", author, text)
	}

	// The same word again counts as agreement, not a second sense.
	again, err := svc.Propose(ctx, contributions.Input{Language: "bem", Headword: "zzamenshi", PartOfSpeech: "noun", GlossEn: "zz water"})
	if err != nil {
		t.Fatal(err)
	}
	if again.NewSense || again.SenseID != r.SenseID {
		t.Errorf("second proposal = %+v, want the existing sense", again)
	}
}

// TestProposeCorrection: a correction is a new candidate linked to the
// meaning it corrects, which stays as it was (US-18).
func TestProposeCorrection(t *testing.T) {
	tx := dbtest.Tx(t)
	ctx := contributor(t, tx)
	svc := contributions.NewService(tx)
	original, err := svc.Propose(ctx, contributions.Input{Language: "bem", Headword: "zzamensi", PartOfSpeech: "noun", GlossEn: "zz water"})
	if err != nil {
		t.Fatal(err)
	}
	fix := func(headword, lang string) contributions.Input {
		return contributions.Input{
			Language: lang, Headword: headword, PartOfSpeech: "noun", GlossEn: "zz water",
			Corrects: uuid.NullUUID{UUID: original.SenseID, Valid: true},
		}
	}
	r, err := svc.Propose(ctx, fix("zzamenshi", "bem"))
	if err != nil || !r.NewSense {
		t.Fatalf("correction = %+v, %v; want a new sense", r, err)
	}
	var corrects uuid.UUID
	var status string
	if err := tx.QueryRow(t.Context(), `SELECT c.corrects, o.status FROM items c, items o WHERE c.id = $1 AND o.id = $2`,
		r.SenseID, original.SenseID).Scan(&corrects, &status); err != nil {
		t.Fatal(err)
	}
	if corrects != original.SenseID || status != "sourced" {
		t.Errorf("corrects %s, original %s; want linked and the original unchanged", corrects, status)
	}
	for name, in := range map[string]contributions.Input{
		"unchanged":      fix("zzamensi", "bem"),
		"other language": fix("zzmadzi", "nya"),
	} {
		if _, err := svc.Propose(ctx, in); httpx.Status(err) != 400 {
			t.Errorf("%s: err = %v, want 400", name, err)
		}
	}
}

func TestProposeRefused(t *testing.T) {
	tx := dbtest.Tx(t)
	svc := contributions.NewService(tx)
	contributor := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "c1", Role: gen.RoleContributor})
	developer := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "d", Role: gen.RoleDeveloper})
	word := func(lang string) contributions.Input {
		return contributions.Input{Language: lang, Headword: "zz", PartOfSpeech: "noun", GlossEn: "zz"}
	}
	tests := []struct {
		name   string
		ctx    context.Context
		in     contributions.Input
		status int
	}{
		{"anonymous", t.Context(), water, 401},
		{"developer", developer, water, 403},
		{"planned language", contributor, word("loz"), 400},
		{"unknown language", contributor, word("zzz"), 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Propose(tt.ctx, tt.in)
			if httpx.Status(err) != tt.status {
				t.Errorf("err = %v, want status %d", err, tt.status)
			}
		})
	}
}
