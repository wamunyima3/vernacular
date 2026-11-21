package accounts_test

import (
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// TestAgreement: u1 agrees with 16 of 21 settled outcomes (76%), under the
// 80% threshold, so their approval is suspended until a senior retrains them.
func TestAgreement(t *testing.T) {
	tx := dbtest.Tx(t)
	agreementFixture(t, tx)
	store := accounts.NewStore(tx)
	svc := accounts.NewService(tx)
	senior := accounts.WithPrincipal(t.Context(), accounts.Principal{
		ID: "u2", Approvals: []accounts.Approval{{Language: "bem", Level: gen.ReviewerLevelSenior}},
	})

	p, err := store.Authenticate(t.Context(), "live")
	if err != nil {
		t.Fatal(err)
	}
	if p.CanReview("bem") || len(p.Approvals) != 1 || !p.Approvals[0].Suspended {
		t.Fatalf("approvals = %+v, want bem suspended", p.Approvals)
	}
	reviewers, err := svc.Reviewers(senior, "bem")
	if err != nil {
		t.Fatal(err)
	}
	if got := reviewers[0]; got.ID != "u1" || got.Decided != 21 || got.Agreed != 16 || !got.Suspended {
		t.Errorf("first reviewer = %+v, want u1 suspended at 16/21", got)
	}

	if err := svc.Retrain(senior, "u2", "bem"); httpx.Status(err) != 403 {
		t.Errorf("retraining yourself: %v, want 403", err)
	}
	if err := svc.Retrain(accounts.WithPrincipal(t.Context(), p), "u1", "bem"); httpx.Status(err) != 403 {
		t.Errorf("a reviewer retrains: %v, want 403", err)
	}
	if err := svc.Retrain(senior, "u1", "bem"); err != nil {
		t.Fatal(err)
	}
	if p, _ = store.Authenticate(t.Context(), "live"); !p.CanReview("bem") {
		t.Errorf("after retraining: %+v, want bem active", p.Approvals)
	}
}

// agreementFixture: reviewer u1 and senior u2 of Bemba; u1 agrees with 16 of
// 21 settled outcomes.
func agreementFixture(t *testing.T, tx pgx.Tx) {
	t.Helper()
	// review inserts n items of a status, each reviewed by u1 with decision.
	review := func(n int, status, decision string) string {
		return fmt.Sprintf(`WITH i AS (INSERT INTO items (kind, language_code, status)
			SELECT 'sense', 'bem', '%s' FROM generate_series(1, %d) RETURNING id)
			INSERT INTO reviews (item_id, reviewer_id, decision, comment, correction, created_at)
			SELECT id, 'u1', '%[3]s', '', CASE WHEN '%[3]s' = 'correct' THEN '{}'::jsonb END, now() - interval '1 hour'
			FROM i`, status, n, decision)
	}
	for _, sql := range []string{
		`INSERT INTO users (id, name, email, email_verified, role) VALUES
			('u1', 'Mwila', 'mwila@test.invalid', true, 'contributor'),
			('u2', 'Senior', 'senior@test.invalid', true, 'contributor')`,
		`INSERT INTO sessions (id, token, user_id, expires_at, updated_at) VALUES ('s1', 'live', 'u1', now() + interval '1 day', now())`,
		`INSERT INTO reviewer_approvals (user_id, language_code, level) VALUES ('u1', 'bem', 'reviewer'), ('u2', 'bem', 'senior')`,
		review(15, "verified", "accept"), // agree
		review(1, "rejected", "correct"), // agree
		review(5, "verified", "reject"),  // disagree
		review(3, "sourced", "reject"),   // not settled
	} {
		if _, err := tx.Exec(t.Context(), sql); err != nil {
			t.Fatal(err)
		}
	}
}
