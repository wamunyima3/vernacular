package review_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/review"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

const quarter = "2099-Q1" // a quarter no real sample uses

func TestQuarter(t *testing.T) {
	for in, want := range map[string]string{"2026-01-01": "2026-Q1", "2026-03-31": "2026-Q1", "2026-04-01": "2026-Q2", "2026-12-31": "2026-Q4"} {
		d, _ := time.Parse(time.DateOnly, in)
		if got := review.Quarter(d); got != want {
			t.Errorf("Quarter(%s) = %s, want %s", in, got, want)
		}
	}
}

// sampled returns a service where water, verified by r1 and r2, has been
// drawn for a re-check, as the only verified item in "tst".
func sampled(t *testing.T) (*review.Service, pgx.Tx, uuid.UUID) {
	tx := dbtest.Tx(t)
	water, _ := fixture(t, tx)
	svc := review.NewService(tx)
	decide(t, svc, "r1", water, accept)
	decide(t, svc, "r2", water, accept)
	if _, err := svc.Sample(reviewer(t), quarter); httpx.Status(err) != 403 {
		t.Fatalf("a reviewer draws a sample: %v, want 403", err)
	}
	admin := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "a1", Role: gen.RoleAdmin})
	if n, err := svc.Sample(admin, quarter); err != nil || n == 0 {
		t.Fatalf("Sample = %d, %v; want at least water", n, err)
	}
	if n, err := svc.Sample(admin, quarter); err != nil || n != 0 {
		t.Fatalf("Sample again = %d, %v; want 0, a quarter is drawn once", n, err)
	}
	return svc, tx, water
}

func TestRecheckQueue(t *testing.T) {
	svc, _, water := sampled(t)
	if _, ok, err := svc.NextRecheck(senior(t, "r1"), "tst", nil); err != nil || ok {
		t.Fatalf("a senior who reviewed water: ok = %v, %v; want nothing", ok, err)
	}
	if _, _, err := svc.NextRecheck(reviewer(t), "tst", nil); httpx.Status(err) != 403 {
		t.Fatalf("a reviewer: %v, want 403", err)
	}
	r, ok, err := svc.NextRecheck(senior(t, "s1"), "tst", nil)
	if err != nil || !ok || r.Candidate.ID != water || r.Quarter != quarter {
		t.Fatalf("NextRecheck = %+v, %v, %v; want water in %s", r, ok, err, quarter)
	}
	if _, ok, _ := svc.NextRecheck(senior(t, "s1"), "tst", []uuid.UUID{water}); ok {
		t.Fatal("a skipped item is offered again")
	}
}

func TestRecordRecheck(t *testing.T) {
	svc, tx, water := sampled(t)
	refused := []struct {
		name    string
		who     string
		item    uuid.UUID
		outcome dbgen.RecheckOutcome
		note    string
		want    int
	}{
		{"a fail without a note", "s1", water, dbgen.RecheckOutcomeFail, " ", 400},
		{"a senior who reviewed it", "r1", water, dbgen.RecheckOutcomePass, "", 403},
		{"an item not drawn", "s1", uuid.New(), dbgen.RecheckOutcomePass, "", 404},
	}
	for _, tt := range refused {
		if _, err := svc.RecordRecheck(senior(t, tt.who), tt.item, tt.outcome, tt.note); httpx.Status(err) != tt.want {
			t.Errorf("%s: %v, want %d", tt.name, err, tt.want)
		}
	}
	c, err := svc.RecordRecheck(senior(t, "s1"), water, dbgen.RecheckOutcomeFail, "Wrong noun class")
	if err != nil || c.Quarter != quarter || c.Status == items.StatusVerified {
		t.Fatalf("a failed re-check = %+v, %v; want water out of verified", c, err)
	}
	if _, err := svc.RecordRecheck(senior(t, "s2"), water, dbgen.RecheckOutcomePass, ""); httpx.Status(err) != 404 {
		t.Fatalf("re-checking twice: %v, want 404", err)
	}
	if _, ok, _ := svc.NextRecheck(senior(t, "s2"), "tst", nil); ok {
		t.Fatal("a checked item is still offered")
	}
	// Back in the queue, only reviews made after the fail count.
	decide(t, svc, "r3", water, accept)
	if item, _, err := items.NewStore(tx).Lock(t.Context(), water); err != nil || item.Status != items.StatusReviewed {
		t.Fatalf("after one new accept: %s, %v; want reviewed", item.Status, err)
	}
}
