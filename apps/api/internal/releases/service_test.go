package releases_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/phrases"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
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
	// Last: a pin's search_path lasts until the test's transaction ends.
	checkPinned(t, tx, svc, first.Name, second.Name)
}

// checkPinned: a pinned read runs the live store on the release's frozen rows
// and names the release in trust; latest is the newest; unknown is not found.
func checkPinned(t *testing.T, tx pgx.Tx, svc *releases.Service, first, second string) {
	t.Helper()
	store := phrases.NewStore(database.Routed{Default: tx})
	pack := func(ctx context.Context) []phrases.Phrase {
		t.Helper()
		got, err := store.Pack(ctx, "bem", "zz-release", "verified")
		if err != nil {
			t.Fatal(err)
		}
		return got
	}
	if live := pack(t.Context()); len(live) != 1 || live[0].Trust.Release != second {
		t.Errorf("live = %+v, want the 1 phrase still verified, in %s", live, second)
	}
	// The first release has both Bemba phrases unedited; the second only the edited one.
	for _, c := range []struct {
		name, want, notes string
		n                 int
	}{{first, first, "", 2}, {"latest", second, "zz edited", 1}} {
		err := svc.Pin(t.Context(), c.name, func(ctx context.Context, release string) error {
			got := pack(ctx)
			if release != c.want || len(got) != c.n || got[0].Trust.Release != c.want || got[0].Notes != c.notes {
				t.Errorf("pinned %s = %s %+v, want %d phrases of %s noted %q", c.name, release, got, c.n, c.want, c.notes)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if err := svc.Pin(t.Context(), "1999.1", nil); httpx.Status(err) != 404 {
		t.Errorf("unknown release: err = %v, want 404", err)
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
