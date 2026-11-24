package releases_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/releases"
)

// TestBuild: a release freezes the verified items only, counts them per language
// and kind, and the next release reports what changed against it.
func TestBuild(t *testing.T) {
	tx := dbtest.Tx(t)
	svc := releases.NewService(tx)
	admin := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "admin-1", Role: gen.RoleAdmin})
	exec := func(sql string) {
		t.Helper()
		if _, err := tx.Exec(t.Context(), sql); err != nil {
			t.Fatal(err)
		}
	}
	// Any verified items already in the database are part of every build; only the
	// changes below are compared.
	exec(`UPDATE items SET status = 'reviewed' WHERE status = 'verified'`)

	if _, err := svc.Build(admin, ""); httpx.Status(err) != 409 {
		t.Fatalf("nothing verified: err = %v, want 409", err)
	}

	exec(`WITH i AS (INSERT INTO items (kind, language_code, status)
			VALUES ('phrase', 'bem', 'verified'), ('phrase', 'bem', 'verified'), ('phrase', 'nya', 'machine')
			RETURNING id, language_code, status)
		INSERT INTO phrases (id, language_code, pack, text, source_en)
		SELECT id, language_code, 'zz-release', 'Zz ' || id, 'Zz ' || status FROM i`)
	first, err := svc.Build(admin, "First release")
	if err != nil {
		t.Fatal(err)
	}
	checkFirst(t, first)

	if _, err := svc.Build(admin, ""); httpx.Status(err) != 409 {
		t.Errorf("no changes: err = %v, want 409", err)
	}

	// One phrase edited, one withdrawn, the Nyanja one verified.
	exec(`UPDATE phrases SET notes = 'zz edited' WHERE id = (
			SELECT id FROM phrases WHERE pack = 'zz-release' AND language_code = 'bem' ORDER BY id LIMIT 1)`)
	exec(`UPDATE items SET status = 'disputed' WHERE id = (
			SELECT id FROM phrases WHERE pack = 'zz-release' AND language_code = 'bem' ORDER BY id DESC LIMIT 1)`)
	exec(`UPDATE items SET status = 'verified' WHERE id IN (
			SELECT id FROM phrases WHERE pack = 'zz-release' AND language_code = 'nya')`)
	second, err := svc.Build(admin, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(second.Name, ".2") || second.Changes != (releases.Changes{Added: 1, Changed: 1, Removed: 1}) {
		t.Errorf("second = %s %+v, want <year>.2 with 1 added, 1 changed, 1 removed", second.Name, second.Changes)
	}

	checkSchema(t, tx, second)

	list, err := svc.List(t.Context())
	if err != nil || len(list) < 2 || list[0].Name != second.Name {
		t.Errorf("list = %v, %v; want newest first", list, err)
	}
}

// checkFirst: the first release has its notes, only the verified items, and
// counts everything as added.
func checkFirst(t *testing.T, first releases.Release) {
	t.Helper()
	if !strings.HasSuffix(first.Name, ".1") || first.Notes != "First release" {
		t.Errorf("first = %s %q, want <year>.1 with the notes", first.Name, first.Notes)
	}
	if got := first.Languages["bem"]["phrase"]; got != 2 || first.Languages["nya"] != nil {
		t.Errorf("counts = %v, want 2 Bemba phrases and no Nyanja", first.Languages)
	}
	if first.Changes != (releases.Changes{Added: 2}) {
		t.Errorf("first changes = %+v, want 2 added", first.Changes)
	}
}

// checkSchema: the release's schema holds its items and the live search indexes (ADR-0012).
func checkSchema(t *testing.T, tx pgx.Tx, r releases.Release) {
	t.Helper()
	schema := releases.Schema(r.Name)
	var phrases, indexes int
	table := pgx.Identifier{schema, "phrases"}.Sanitize() //nolint:misspell // pgx's method name
	err := tx.QueryRow(t.Context(), `SELECT
			(SELECT count(*) FROM `+table+` WHERE pack = 'zz-release'),
			(SELECT count(*) FROM pg_indexes WHERE schemaname = $1 AND tablename = 'entries')`,
		schema).Scan(&phrases, &indexes)
	if err != nil {
		t.Fatal(err)
	}
	if phrases != 2 || indexes == 0 {
		t.Errorf("%s: %d phrases, %d entry indexes; want the 2 verified phrases and the live indexes", schema, phrases, indexes)
	}
}

func TestBuildRequiresAdmin(t *testing.T) {
	svc := releases.NewService(dbtest.Tx(t))
	developer := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "d", Role: gen.RoleDeveloper})
	if _, err := svc.Build(developer, ""); !errors.Is(err, httpx.ErrForbidden) {
		t.Errorf("developer: err = %v, want forbidden", err)
	}
}
