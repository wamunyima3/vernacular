package accounts_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

// qualificationFixture adds language "tst" with n verified words, a taker t1
// and a senior s1.
func qualificationFixture(t *testing.T, tx pgx.Tx, n int) (taker, senior context.Context) {
	t.Helper()
	ctx := t.Context()
	if _, err := tx.Exec(ctx, `
		INSERT INTO languages (code, name, own_name, wave, status) VALUES ('tst', 'Test', 'Test', 9, 'active');
		INSERT INTO users (id, name, email, email_verified, role) VALUES
			('t1', 'Taker', 't1@test.invalid', true, 'contributor'), ('s1', 'Senior', 's1@test.invalid', true, 'contributor');`); err != nil {
		t.Fatal(err)
	}
	source, err := sources.Lookup(ctx, tx, "Community contributors")
	if err != nil {
		t.Fatal(err)
	}
	lex := lexicon.NewStore(tx)
	for i := range n {
		if _, err := lex.Add(ctx, lexicon.Proposal{
			Language: "tst", Headword: fmt.Sprintf("zzword%d", i), PartOfSpeech: "noun",
			GlossEn: fmt.Sprintf("zz meaning %d", i), Origin: items.Contributed(source, "s1"),
		}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := tx.Exec(ctx, `UPDATE items SET status = 'verified' WHERE language_code = 'tst'`); err != nil {
		t.Fatal(err)
	}
	return accounts.WithPrincipal(ctx, accounts.Principal{ID: "t1", Role: gen.RoleContributor}),
		accounts.WithPrincipal(ctx, accounts.Principal{ID: "s1", Role: gen.RoleContributor, Approvals: []accounts.Approval{{Language: "tst", Level: gen.ReviewerLevelSenior}}})
}

// rightAnswers picks, for each question, the fixture word with its meaning.
func rightAnswers(test accounts.Test) []uuid.UUID {
	answers := make([]uuid.UUID, len(test.Questions))
	for i, q := range test.Questions {
		for _, o := range q.Options {
			var n int
			_, _ = fmt.Sscanf(o.Headword, "zzword%d", &n)
			if q.GlossEn == fmt.Sprintf("zz meaning %d", n) {
				answers[i] = o.ID
			}
		}
	}
	return answers
}

// TestQualification: a taker who passes is approved once a senior confirms
// (DQ-3, US-24).
func TestQualification(t *testing.T) {
	tx := dbtest.Tx(t)
	taker, senior := qualificationFixture(t, tx, 20)
	svc := accounts.NewService(tx)

	test, err := svc.StartTest(taker, "tst")
	if err != nil {
		t.Fatal(err)
	}
	if len(test.Questions) != 10 || len(test.Questions[0].Options) != 4 {
		t.Fatalf("test = %+v, want 10 questions of 4 options", test)
	}
	answers := rightAnswers(test)
	answers[0] = uuid.New() // one wrong still passes
	r, err := svc.AnswerTest(taker, test.ID, answers)
	if err != nil || r.Score != 9 || !r.Passed {
		t.Fatalf("result = %+v, %v; want 9 of 10, passed", r, err)
	}
	if _, err := svc.AnswerTest(taker, test.ID, answers); httpx.Status(err) != 409 {
		t.Errorf("second hand-in: %v, want 409", err)
	}
	if _, err := svc.StartTest(taker, "tst"); httpx.Status(err) != 409 {
		t.Errorf("start while a pass waits: %v, want 409", err)
	}

	checkConfirm(taker, senior, t, tx, svc, test.ID)
}

func TestQualificationRefused(t *testing.T) {
	tx := dbtest.Tx(t)
	taker, _ := qualificationFixture(t, tx, 12)
	svc := accounts.NewService(tx)
	if _, err := svc.StartTest(taker, "tst"); httpx.Status(err) != 409 {
		t.Errorf("12 words: %v, want 409 (not enough for a test)", err)
	}
	if _, err := tx.Exec(t.Context(), `INSERT INTO qualification_attempts (user_id, language_code, expected, score, passed, answered_at)
		VALUES ('t1', 'tst', '{}', 2, false, now())`); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.StartTest(taker, "tst"); httpx.Status(err) != 409 {
		t.Errorf("retake the same day: %v, want 409", err)
	}
}

// checkConfirm: a senior, not the taker, confirms a pass, approving the taker.
func checkConfirm(taker, senior context.Context, t *testing.T, tx pgx.Tx, svc *accounts.Service, id uuid.UUID) {
	t.Helper()
	pending, err := svc.PendingTests(senior)
	if err != nil || len(pending) != 1 || pending[0].UserID != "t1" {
		t.Fatalf("pending = %+v, %v; want t1's pass", pending, err)
	}
	if err := svc.ConfirmTest(taker, id); httpx.Status(err) != 403 {
		t.Errorf("taker confirms: %v, want 403", err)
	}
	if err := svc.ConfirmTest(senior, id); err != nil {
		t.Fatal(err)
	}
	var level, by string
	if err := tx.QueryRow(t.Context(), `SELECT level, approved_by FROM reviewer_approvals WHERE user_id = 't1' AND language_code = 'tst'`).
		Scan(&level, &by); err != nil || level != "reviewer" || by != "s1" {
		t.Errorf("approval = %s by %s, %v; want reviewer by s1", level, by, err)
	}
}
