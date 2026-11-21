package review_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/review"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

func senior(t *testing.T, id string) context.Context {
	return accounts.WithPrincipal(t.Context(), accounts.Principal{
		ID: id, Role: gen.RoleContributor, Approvals: []accounts.Approval{{Language: "tst", Level: gen.ReviewerLevelSenior}},
	})
}

// disputed returns a service where r1 accepted and r2 rejected both water
// and the first stone candidate, so both are disputed.
func disputed(t *testing.T) (svc *review.Service, status func(uuid.UUID) items.Status, water, stone uuid.UUID) {
	tx := dbtest.Tx(t)
	water, _ = fixture(t, tx)
	stone = stoneCard(t, tx)[0]
	svc = review.NewService(tx)
	for _, id := range []uuid.UUID{water, stone} {
		decide(t, svc, "r1", id, accept)
		decide(t, svc, "r2", id, reject)
	}
	status = func(id uuid.UUID) items.Status {
		t.Helper()
		item, _, err := items.NewStore(tx).Lock(t.Context(), id)
		if err != nil {
			t.Fatal(err)
		}
		return item.Status
	}
	return svc, status, water, stone
}

func decide(t *testing.T, svc *review.Service, who string, id uuid.UUID, d dbgen.ReviewDecision) {
	t.Helper()
	if _, err := svc.Decide(reviewer(t, who), []review.Decision{{SenseID: id, Decision: d}}, "why", nil); err != nil {
		t.Fatal(err)
	}
}

func TestDisputes(t *testing.T) {
	svc, _, water, _ := disputed(t)
	disputes, err := svc.Disputes(senior(t, "s1"))
	if err != nil {
		t.Fatal(err)
	}
	if len(disputes) != 2 || len(disputes[0].Reviews) != 2 || disputes[0].Candidate.GlossEn == "" {
		t.Fatalf("disputes = %+v, want 2 with 2 reviews and a gloss", disputes)
	}
	tests := []struct {
		name string
		err  func() error
		want int
	}{
		{"a reviewer lists disputes", func() error { _, err := svc.Disputes(reviewer(t)); return err }, 403},
		{"a senior who reviewed it", func() error { _, err := svc.Resolve(senior(t, "r1"), water, keep, "fine"); return err }, 403},
		{"an empty note", func() error { _, err := svc.Resolve(senior(t, "s1"), water, keep, "  "); return err }, 400},
		{"not a sense", func() error { _, err := svc.Resolve(senior(t, "s1"), uuid.New(), keep, "fine"); return err }, 404},
	}
	for _, tt := range tests {
		if err := tt.err(); httpx.Status(err) != tt.want {
			t.Errorf("%s: %v, want %d", tt.name, err, tt.want)
		}
	}
}

func TestKeep(t *testing.T) {
	svc, status, water, _ := disputed(t)
	if got, err := svc.Resolve(senior(t, "s1"), water, keep, "Common word"); err != nil || got != items.StatusVerified {
		t.Fatalf("keep = %s, %v; want verified", got, err)
	}
	if _, err := svc.Resolve(senior(t, "s2"), water, sendBack, "again"); httpx.Status(err) != 409 {
		t.Errorf("resolving a settled item: %v, want 409", err)
	}
	if status(water) != items.StatusVerified {
		t.Errorf("after keep: %s, want verified", status(water))
	}
}

func TestSendBack(t *testing.T) {
	svc, status, _, stone := disputed(t)
	if got, err := svc.Resolve(senior(t, "s1"), stone, sendBack, "Ask again"); err != nil || got != items.StatusMachine {
		t.Fatalf("send back = %s, %v; want its origin, machine", got, err)
	}
	decide(t, svc, "r3", stone, accept)
	if got := status(stone); got != items.StatusReviewed {
		t.Errorf("one accept after a send-back: %s, want reviewed (earlier reviews set aside)", got)
	}
	if disputes, _ := svc.Disputes(senior(t, "s1")); len(disputes) != 1 {
		t.Errorf("disputes left: %d, want 1 (water)", len(disputes))
	}
}
