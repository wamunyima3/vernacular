package review_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/review"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

// fixture builds a queue in a test language "tst":
//
//	zz water: Claude and Google agree on amenshi       → quick
//	zz fire:  Claude only, umulilo                     → full
//	zz stone: Claude ibwe, Google libwe                → full, two candidates
//	zz tree:  icimuti, contributed by the reviewer     → never shown to them
func fixture(t *testing.T, tx pgx.Tx) (water uuid.UUID, tree uuid.UUID) {
	t.Helper()
	ctx := t.Context()
	if _, err := tx.Exec(ctx, `
		INSERT INTO languages (code, name, own_name, wave, status) VALUES ('tst', 'Test', 'Test', 9, 'active');
		INSERT INTO users (id, name, email, email_verified, role) VALUES ('r1', 'Reviewer', 'r1@test.invalid', true, 'contributor'),
			('r2', 'Second', 'r2@test.invalid', true, 'contributor'), ('r3', 'Third', 'r3@test.invalid', true, 'contributor');`); err != nil {
		t.Fatal(err)
	}
	source := func(name string) uuid.UUID {
		id, err := sources.Lookup(ctx, tx, name)
		if err != nil {
			t.Fatal(err)
		}
		return id
	}
	claude, google, people := source("Claude"), source("Google Translate"), source("Community contributors")
	lex := lexicon.NewStore(tx)
	add := func(headword, gloss string, origin items.Origin) uuid.UUID {
		added, err := lex.Add(ctx, lexicon.Proposal{Language: "tst", Headword: headword, PartOfSpeech: "noun", GlossEn: gloss, Origin: origin})
		if err != nil {
			t.Fatal(err)
		}
		return added.SenseID
	}
	water = add("amenshi", "zz water", items.Machine(claude, "m", "c1"))
	add("amenshi", "zz water", items.Machine(google, "g", "g1"))
	add("umulilo", "zz fire", items.Machine(claude, "m", "c2"))
	add("ibwe", "zz stone", items.Machine(claude, "m", "c3"))
	add("libwe", "zz stone", items.Machine(google, "g", "g3"))
	tree = add("icimuti", "zz tree", items.Contributed(people, "r1"))
	return water, tree
}

// reviewer signs in an approved reviewer of "tst"; the fixture's user is r1.
func reviewer(t *testing.T, id ...string) context.Context {
	return accounts.WithPrincipal(t.Context(), accounts.Principal{
		ID: append(id, "r1")[0], Role: gen.RoleContributor, Approvals: []accounts.Approval{{Language: "tst", Level: gen.ReviewerLevelReviewer}},
	})
}

func TestQueue(t *testing.T) {
	tx := dbtest.Tx(t)
	water, _ := fixture(t, tx)
	svc := review.NewService(tx)
	ctx := reviewer(t)

	card, ok, err := svc.Next(ctx, "tst", nil)
	if err != nil || !ok {
		t.Fatalf("Next = %v, %v", ok, err)
	}
	if card.GlossEn != "zz water" || card.Priority != dbgen.ReviewPriorityQuick || len(card.Candidates) != 1 || card.Candidates[0].ID != water {
		t.Fatalf("first card = %+v, want the quick confirm for zz water", card)
	}
	if _, err := svc.Decide(ctx, []review.Decision{{ItemID: water, Decision: dbgen.ReviewDecisionAccept}}, "", nil); err != nil {
		t.Fatal(err)
	}

	// The rest are full reviews; skipping each in turn walks the queue.
	var skip []uuid.UUID
	seen := map[string]int{}
	for {
		card, ok, err := svc.Next(ctx, "tst", skip)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			break
		}
		if card.Priority != dbgen.ReviewPriorityFull {
			t.Errorf("%s priority = %s, want full", card.GlossEn, card.Priority)
		}
		seen[card.GlossEn] = len(card.Candidates)
		skip = append(skip, card.ConceptID)
	}
	if want := map[string]int{"zz fire": 1, "zz stone": 2}; len(seen) != len(want) || seen["zz fire"] != 1 || seen["zz stone"] != 2 {
		t.Errorf("cards after water = %v, want %v (never water again, never the reviewer's own tree)", seen, want)
	}
}

func TestDecideRefused(t *testing.T) {
	tx := dbtest.Tx(t)
	water, tree := fixture(t, tx)
	svc := review.NewService(tx)
	ctx := reviewer(t)
	outsider := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "r2", Role: gen.RoleContributor})

	if _, _, err := svc.Next(outsider, "tst", nil); !errors.Is(err, httpx.ErrForbidden) {
		t.Errorf("Next without approval = %v, want forbidden", err)
	}
	tests := []struct {
		name   string
		ctx    context.Context
		id     uuid.UUID
		status int
	}{
		{"not approved for the language", outsider, water, 403},
		{"own contribution", ctx, tree, 403},
		{"not a sense", ctx, uuid.New(), 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Decide(tt.ctx, []review.Decision{{ItemID: tt.id, Decision: dbgen.ReviewDecisionAccept}}, "", nil)
			if httpx.Status(err) != tt.status {
				t.Errorf("err = %v, want status %d", err, tt.status)
			}
		})
	}
}

func TestDecideTwice(t *testing.T) {
	tx := dbtest.Tx(t)
	water, _ := fixture(t, tx)
	svc := review.NewService(tx)
	if _, err := svc.Decide(reviewer(t), []review.Decision{{ItemID: water, Decision: dbgen.ReviewDecisionAccept}}, "", nil); err != nil {
		t.Fatal(err)
	}
	_, err := svc.Decide(reviewer(t), []review.Decision{{ItemID: water, Decision: dbgen.ReviewDecisionReject}}, "wrong", nil)
	if httpx.Status(err) != 409 {
		t.Errorf("second decision = %v, want conflict", err)
	}
}

func TestUndo(t *testing.T) {
	tx := dbtest.Tx(t)
	water, _ := fixture(t, tx)
	svc := review.NewService(tx)
	decide := func() uuid.UUID {
		t.Helper()
		r, err := svc.Decide(reviewer(t), []review.Decision{{ItemID: water, Decision: dbgen.ReviewDecisionAccept}}, "", nil)
		if err != nil {
			t.Fatal(err)
		}
		return r[0].ID
	}

	id := decide()
	other := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "r2", Role: gen.RoleContributor})
	if got := httpx.Status(svc.Undo(other, []uuid.UUID{id})); got != 404 {
		t.Errorf("undo by another reviewer: status %d, want 404", got)
	}
	if err := svc.Undo(reviewer(t), []uuid.UUID{id}); err != nil {
		t.Fatal(err)
	}
	if card, ok, _ := svc.Next(reviewer(t), "tst", nil); !ok || card.GlossEn != "zz water" {
		t.Errorf("after undo the next card is %q, want zz water back", card.GlossEn)
	}

	id = decide()
	if _, err := tx.Exec(t.Context(), `UPDATE reviews SET created_at = now() - interval '11 minutes' WHERE id = $1`, id); err != nil {
		t.Fatal(err)
	}
	if got := httpx.Status(svc.Undo(reviewer(t), []uuid.UUID{id})); got != 409 {
		t.Errorf("undo after the window: status %d, want 409", got)
	}
}
