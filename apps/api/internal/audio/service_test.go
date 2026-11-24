package audio_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/audio"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/blob"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

const pdf = "%PDF-1.4 signed consent form"

// TestConsent: an admin records a speaker's signed form, the speaker sees it
// and withdraws, and a new form restores consent (NFR-6, US-33). With audio
// off, nothing answers.
func TestConsent(t *testing.T) {
	tx := dbtest.Tx(t)
	if _, err := tx.Exec(t.Context(), `INSERT INTO users (id, name, email, email_verified, role) VALUES
		('admin', 'A', 'a@test.invalid', true, 'admin'), ('sp', 'S', 's@test.invalid', true, 'contributor')`); err != nil {
		t.Fatal(err)
	}
	bucket, err := blob.Open(t.Context(), "mem://")
	if err != nil {
		t.Fatal(err)
	}
	admin := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "admin", Role: gen.RoleAdmin})
	sp := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "sp", Role: gen.RoleContributor})

	if _, err := audio.NewService(tx, bucket, false).Me(sp); !errors.Is(err, audio.ErrDisabled) {
		t.Fatalf("audio off: %v, want ErrDisabled", err)
	}
	svc := audio.NewService(tx, bucket, true)
	checkRefusals(admin, sp, t, svc)

	if _, err := svc.RecordConsent(admin, "sp", strings.NewReader(pdf)); err != nil {
		t.Fatal(err)
	}
	if me, err := svc.Me(sp); err != nil || me.WithdrawnAt != nil {
		t.Fatalf("me = %+v, %v; want consent in force", me, err)
	}
	rd, _, err := svc.ConsentForm(admin, "sp")
	if err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(rd)
	_ = rd.Close()
	if string(got) != pdf {
		t.Errorf("form = %q", got)
	}

	if w, err := svc.Withdraw(sp); err != nil || w.WithdrawnAt == nil {
		t.Fatalf("withdraw = %+v, %v", w, err)
	}
	if _, err := svc.Withdraw(sp); httpx.Status(err) != 409 {
		t.Errorf("withdraw twice: %v, want 409", err)
	}
	if again, err := svc.RecordConsent(admin, "sp", strings.NewReader(pdf)); err != nil || again.WithdrawnAt != nil {
		t.Errorf("new form after withdrawal = %+v, %v; want consent restored", again, err)
	}
}

// checkRefusals: only admins record consent, for real users, from a PDF or
// image; before that the speaker has no record.
func checkRefusals(admin, sp context.Context, t *testing.T, svc *audio.Service) {
	t.Helper()
	for _, tc := range []struct {
		name   string
		ctx    context.Context
		user   string
		form   string
		status int
	}{
		{"speaker records own consent", sp, "sp", pdf, 403},
		{"unknown user", admin, "nobody", pdf, 404},
		{"text form", admin, "sp", "plain text", 400},
	} {
		if _, err := svc.RecordConsent(tc.ctx, tc.user, strings.NewReader(tc.form)); httpx.Status(err) != tc.status {
			t.Errorf("%s: %v, want %d", tc.name, err, tc.status)
		}
	}
	if _, err := svc.Me(sp); httpx.Status(err) != 404 {
		t.Errorf("before consent: %v, want 404", err)
	}
}
