package releases_test

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/releases"
)

// TestExport: the download holds the release's frozen rows as CSV, with a
// README and release.json (NFR-9, US-17).
func TestExport(t *testing.T) {
	tx := dbtest.Tx(t)
	svc := releases.NewService(tx)
	if _, err := tx.Exec(t.Context(), `UPDATE items SET status = 'reviewed' WHERE status = 'verified';
		WITH i AS (INSERT INTO items (kind, language_code, status) VALUES ('phrase', 'bem', 'verified') RETURNING id)
		INSERT INTO phrases (id, language_code, pack, text, source_en) SELECT id, 'bem', 'zz-export', 'Mwapoleni, "mukwai"', 'Hello' FROM i`); err != nil {
		t.Fatal(err)
	}
	built, err := svc.Build(accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "a", Role: gen.RoleAdmin}), "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(t.Context(), `UPDATE phrases SET text = 'Changed later' WHERE pack = 'zz-export'`); err != nil {
		t.Fatal(err)
	}

	name, body, err := svc.Export(t.Context(), "latest")
	if err != nil || name != built.Name {
		t.Fatalf("export latest = %s, %v; want %s", name, err, built.Name)
	}
	files := unzip(t, body)
	if want := `"Mwapoleni, ""mukwai"""`; !strings.Contains(files["phrases.csv"], want) {
		t.Errorf("phrases.csv = %q, want the frozen text quoted as %s", files["phrases.csv"], want)
	}
	if !strings.HasPrefix(files["entries.csv"], "id,language,headword,") {
		t.Errorf("entries.csv starts %q, want its header", files["entries.csv"])
	}
	if !strings.Contains(files["README.md"], built.Name) || !strings.Contains(files["release.json"], `"name": "`+built.Name) {
		t.Errorf("README or release.json does not name %s", built.Name)
	}

	if _, _, err := svc.Export(t.Context(), "1999.1"); httpx.Status(err) != 404 {
		t.Errorf("unknown release: err = %v, want 404", err)
	}
}

func unzip(t *testing.T, body []byte) map[string]string {
	t.Helper()
	z, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, f := range z.File {
		r, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		b, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
		out[f.Name] = string(b)
	}
	return out
}
