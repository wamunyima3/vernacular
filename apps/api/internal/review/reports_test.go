package review_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/review"
)

// reported: a signed-out user and r2 report water.
func reported(t *testing.T) (svc *review.Service, status func() items.Status, water, anon, own uuid.UUID) {
	t.Helper()
	tx := dbtest.Tx(t)
	water, _ = fixture(t, tx)
	svc = review.NewService(tx)
	status = func() items.Status {
		t.Helper()
		item, _, err := items.NewStore(tx).Lock(t.Context(), water)
		if err != nil {
			t.Fatal(err)
		}
		return item.Status
	}
	var err error
	if anon, err = svc.Report(t.Context(), water, "Means moon"); err != nil {
		t.Fatal(err)
	}
	if own, err = svc.Report(reviewer(t, "r2"), water, "Wrong spelling"); err != nil {
		t.Fatal(err)
	}
	return svc, status, water, anon, own
}

func TestReportInbox(t *testing.T) {
	svc, _, water, _, _ := reported(t)
	if _, err := svc.Report(t.Context(), uuid.New(), "x"); httpx.Status(err) != 404 {
		t.Errorf("reporting nothing: %v, want 404", err)
	}
	if _, err := svc.Reports(accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "r3"}), "tst"); httpx.Status(err) != 403 {
		t.Errorf("inbox for a non-reviewer: %v, want 403", err)
	}
	if inbox, err := svc.Reports(reviewer(t), "tst"); err != nil || len(inbox) != 2 {
		t.Errorf("inbox = %d, %v; want 2", len(inbox), err)
	}
	// A report sent with r1's API key (say, from the portal) is anonymous, so r1 may decide it.
	viaKey, err := svc.Report(accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "r1", KeyID: uuid.New()}), water, "via an app")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.DecideReport(reviewer(t), viaKey, false, "fine"); err != nil {
		t.Errorf("r1 deciding a report sent with their key: %v", err)
	}
}

// TestReportDecisions: r2 cannot decide their own report; r1 dismisses it
// and confirms the other, so water is disputed until a senior keeps it.
func TestReportDecisions(t *testing.T) {
	svc, status, water, anon, own := reported(t)
	if _, err := svc.DecideReport(reviewer(t, "r2"), own, true, "yes"); httpx.Status(err) != 403 {
		t.Errorf("deciding your own report: %v, want 403", err)
	}
	if _, err := svc.DecideReport(reviewer(t), own, false, "Spelling is fine"); err != nil {
		t.Fatal(err)
	}
	if got, err := svc.DecideReport(reviewer(t), anon, true, "Checked: it is moon"); err != nil || got != items.StatusDisputed {
		t.Fatalf("confirm = %s, %v; want disputed", got, err)
	}
	if _, err := svc.DecideReport(reviewer(t), anon, false, "again"); httpx.Status(err) != 409 {
		t.Errorf("deciding twice: %v, want 409", err)
	}
	if disputes, _ := svc.Disputes(senior(t, "s1")); len(disputes) != 1 || disputes[0].Candidate.ID != water {
		t.Errorf("disputes = %+v, want water", disputes)
	}
	if got, err := svc.Resolve(senior(t, "s1"), water, keep, "It is water"); err != nil || got != items.StatusVerified {
		t.Fatalf("keep = %s, %v; want verified (the report is settled)", got, err)
	}
	if got := status(); got != items.StatusVerified {
		t.Errorf("after keep: %s", got)
	}
}
