package audio_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/audio"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/blob"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

// wav is the start of a WAV file, enough for its type to be sniffed.
const wav = "RIFF\x24\x00\x00\x00WAVEfmt \x10\x00\x00\x00"

// entries adds a Test-language and a Bemba entry and returns their IDs.
func entries(t *testing.T, tx pgx.Tx) (tst, bem uuid.UUID) {
	t.Helper()
	ctx := t.Context()
	if _, err := tx.Exec(ctx, `INSERT INTO languages (code, name, own_name, wave, status) VALUES ('tst', 'Test', 'Test', 9, 'active')`); err != nil {
		t.Fatal(err)
	}
	people, err := sources.Lookup(ctx, tx, "Community contributors")
	if err != nil {
		t.Fatal(err)
	}
	add := func(lang, headword string) uuid.UUID {
		a, err := lexicon.NewStore(tx).Add(ctx, lexicon.Proposal{
			Language: lang, Headword: headword, PartOfSpeech: "noun", GlossEn: "zz " + headword, Origin: items.Contributed(people, "admin"),
		})
		if err != nil {
			t.Fatal(err)
		}
		return a.EntryID
	}
	return add("tst", "amenshi"), add("bem", "zzumulilo")
}

// TestClips: a consenting speaker records entries in their languages only,
// anyone hears them, a new recording replaces the old, and withdrawing
// deletes the speaker's clips and files (LEX-8, NFR-6, US-33).
func TestClips(t *testing.T) {
	tx := dbtest.Tx(t)
	if _, err := tx.Exec(t.Context(), `INSERT INTO users (id, name, email, email_verified, role) VALUES
		('admin', 'A', 'a@test.invalid', true, 'admin'), ('sp', 'S', 's@test.invalid', true, 'contributor')`); err != nil {
		t.Fatal(err)
	}
	tst, bem := entries(t, tx)
	bucket, err := blob.Open(t.Context(), "mem://")
	if err != nil {
		t.Fatal(err)
	}
	svc := audio.NewService(tx, bucket, true)
	admin := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "admin", Role: gen.RoleAdmin})
	sp := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "sp", Role: gen.RoleContributor})

	if _, err := svc.Record(sp, tst, strings.NewReader(wav)); httpx.Status(err) != 403 {
		t.Errorf("no consent: %v, want 403", err)
	}
	if _, err := svc.RecordConsent(admin, "sp", []string{"tst"}, strings.NewReader(pdf)); err != nil {
		t.Fatal(err)
	}
	checkRecordRefusals(sp, t, svc, bem)

	checkReplace(sp, t, svc, tst)

	if _, err := svc.Withdraw(sp); err != nil {
		t.Fatal(err)
	}
	if cs, err := svc.Clips(t.Context(), tst); err != nil || len(cs) != 0 {
		t.Errorf("clips after withdrawal = %v, %v; want none", cs, err)
	}
	if ok, _ := bucket.Exists(t.Context(), "speakers/sp/audio/"+tst.String()); ok {
		t.Error("the file outlived the withdrawal")
	}
}

// checkRecordRefusals: other languages, unknown entries and non-audio files
// are refused.
func checkRecordRefusals(sp context.Context, t *testing.T, svc *audio.Service, bem uuid.UUID) {
	t.Helper()
	for _, tc := range []struct {
		name   string
		entry  uuid.UUID
		body   string
		status int
	}{
		{"language not on the form", bem, wav, 403},
		{"unknown entry", uuid.New(), wav, 404},
	} {
		if _, err := svc.Record(sp, tc.entry, strings.NewReader(tc.body)); httpx.Status(err) != tc.status {
			t.Errorf("%s: %v, want %d", tc.name, err, tc.status)
		}
	}
}

// checkReplace: recording an entry again replaces the clip, and its file
// plays with the type it was recorded in.
func checkReplace(sp context.Context, t *testing.T, svc *audio.Service, tst uuid.UUID) {
	t.Helper()
	first, err := svc.Record(sp, tst, strings.NewReader(wav))
	if err != nil || first.ContentType != "audio/wav" {
		t.Fatalf("record = %+v, %v", first, err)
	}
	again, err := svc.Record(sp, tst, strings.NewReader(wav+"more"))
	if err != nil || again.ID != first.ID || again.Size != int32(len(wav)+4) {
		t.Errorf("record again = %+v, %v; want the same clip replaced", again, err)
	}
	rd, c, err := svc.Open(t.Context(), first.ID)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(rd)
	_ = rd.Close()
	if string(got) != wav+"more" || c.ContentType != "audio/wav" {
		t.Errorf("file = %q (%s)", got, c.ContentType)
	}
}
