package review_test

import (
	"context"
	"slices"
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
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
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

func admin(t *testing.T) context.Context {
	return accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "a1", Role: gen.RoleAdmin})
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
	admin := admin(t)
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

// TestSendBackBatch fails water, from Claude model "m", then sends the rest
// of that batch back: fire, verified after the sample was drawn.
func TestSendBackBatch(t *testing.T) {
	svc, tx, water := sampled(t)
	var fire uuid.UUID
	if err := tx.QueryRow(t.Context(), `SELECT s.item_id FROM item_sources s JOIN items i ON i.id = s.item_id WHERE s.reference = 'c2' AND i.kind = 'sense'`).Scan(&fire); err != nil {
		t.Fatal(err)
	}
	decide(t, svc, "r1", fire, accept)
	decide(t, svc, "r2", fire, accept)
	claude, err := sources.Lookup(t.Context(), tx, "Claude")
	if err != nil {
		t.Fatal(err)
	}
	batch := review.Batch{Language: "tst", SourceID: claude, Model: "m"}
	if _, err := svc.SendBackBatch(admin(t), quarter, batch); httpx.Status(err) != 409 {
		t.Fatalf("a batch with no failed checks: %v, want 409", err)
	}
	if _, err := svc.RecordRecheck(senior(t, "s1"), water, dbgen.RecheckOutcomeFail, "Wrong"); err != nil {
		t.Fatal(err)
	}
	checkBatches(t, svc, claude)
	if _, err := svc.SendBackBatch(senior(t, "s1"), quarter, batch); httpx.Status(err) != 403 {
		t.Fatalf("a senior reviewer sends back: %v, want 403", err)
	}
	if n, err := svc.SendBackBatch(admin(t), quarter, batch); err != nil || n != 1 {
		t.Fatalf("SendBackBatch = %d, %v; want fire", n, err)
	}
	if item, _, err := items.NewStore(tx).Lock(t.Context(), fire); err != nil || item.Status == items.StatusVerified {
		t.Fatalf("fire is %s, %v; want it back in review", item.Status, err)
	}
	if _, err := svc.SendBackBatch(admin(t), quarter, batch); httpx.Status(err) != 409 {
		t.Fatalf("sending back twice: %v, want 409", err)
	}
}

func checkBatches(t *testing.T, svc *review.Service, claude uuid.UUID) {
	t.Helper()
	rows, err := svc.Batches(admin(t), quarter)
	if err != nil {
		t.Fatal(err)
	}
	i := slices.IndexFunc(rows, func(r dbgen.RecheckBatchesRow) bool { return r.LanguageCode == "tst" })
	if i < 0 || rows[i].SourceID != claude || rows[i].Model != "m" || rows[i].Checked != 1 || rows[i].Failed != 1 {
		t.Fatalf("batches = %+v, want tst Claude m with 1 of 1 failed", rows)
	}
	if !review.Fails(1, 1) || review.Fails(100, 2) || !review.Fails(100, 3) {
		t.Fatal("Fails should mean more than 2%")
	}
}
