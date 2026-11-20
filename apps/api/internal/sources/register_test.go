package sources_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/blob"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

var input = sources.Input{
	Kind: "licensed", Name: "Zz Test Dictionary", RightsHolder: "Author", Licence: "CC BY 4.0",
	AllowedUses: []string{"redistribute"},
}

// pdf is enough of a PDF for content sniffing.
var pdf = []byte("%PDF-1.7\n1 0 obj\n<<>>\nendobj\n")

func newRegister(t *testing.T) *sources.Register {
	t.Helper()
	bucket, err := blob.Open(t.Context(), "mem://")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = bucket.Close() })
	return sources.NewRegister(dbtest.Tx(t), bucket)
}

func TestRegisterRequiresAdmin(t *testing.T) {
	reg := newRegister(t)
	developer := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "d", Role: gen.RoleDeveloper})
	if _, err := reg.Create(t.Context(), input); !errors.Is(err, httpx.ErrUnauthorized) {
		t.Errorf("anonymous create: err = %v, want unauthorized", err)
	}
	if _, err := reg.Create(developer, input); !errors.Is(err, httpx.ErrForbidden) {
		t.Errorf("developer create: err = %v, want forbidden", err)
	}
	if _, err := reg.List(t.Context()); err != nil {
		t.Errorf("anonymous list: %v, want the public register", err)
	}
}

func TestRegisterLifecycle(t *testing.T) {
	reg := newRegister(t)
	ctx := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "admin-1", Role: gen.RoleAdmin})

	src, err := reg.Create(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := reg.Create(ctx, input); httpx.Status(err) != 409 {
		t.Errorf("duplicate name: err = %v, want 409", err)
	}

	edited := input
	edited.Licence = "CC BY-SA 4.0"
	if got, err := reg.Update(ctx, src.ID, edited); err != nil || got.Licence != "CC BY-SA 4.0" {
		t.Errorf("update = %+v, %v", got, err)
	}
	if _, err := reg.Update(ctx, uuid.New(), edited); !errors.Is(err, httpx.ErrNotFound) {
		t.Errorf("update missing: err = %v, want not found", err)
	}

	uploads := []struct {
		name       string
		body       []byte
		wantStatus int
	}{
		{"not a document", []byte("hello"), 400},
		{"too large", append(bytes.Clone(pdf), make([]byte, sources.MaxPermissionSize)...), 413},
		{"pdf", pdf, 0},
	}
	for _, u := range uploads {
		_, err := reg.UploadPermission(ctx, src.ID, bytes.NewReader(u.body))
		if httpx.Status(err) != u.wantStatus || (u.wantStatus == 0 && err != nil) {
			t.Errorf("upload %s: err = %v, want status %d", u.name, err, u.wantStatus)
		}
	}
	rd, _, err := reg.Permission(ctx, src.ID)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(rd)
	_ = rd.Close()
	if !bytes.Equal(got, pdf) {
		t.Errorf("downloaded %q, want the uploaded PDF", got)
	}
	if s, _ := reg.Get(ctx, src.ID); !s.HasPermissionDocument {
		t.Error("has_permission_document = false after upload")
	}

	if err := reg.Delete(ctx, src.ID); err != nil {
		t.Errorf("delete unused source: %v", err)
	}
	if err := reg.Delete(ctx, src.ID); !errors.Is(err, httpx.ErrNotFound) {
		t.Errorf("delete again: err = %v, want not found", err)
	}
}
