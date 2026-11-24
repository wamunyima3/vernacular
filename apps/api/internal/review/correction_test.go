package review_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/review"
)

// stoneCard returns the senses of "zz stone", where Claude and Google disagree.
func stoneCard(t *testing.T, tx pgx.Tx) []uuid.UUID {
	t.Helper()
	rows, err := tx.Query(t.Context(), `SELECT s.id FROM senses s JOIN concepts c ON c.id = s.concept_id
		WHERE c.gloss_en = 'zz stone' ORDER BY s.id`)
	if err != nil {
		t.Fatal(err)
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		t.Fatal(err)
	}
	return ids
}

// describe reads a sense as "headword noun-class by author".
func describe(t *testing.T, tx pgx.Tx, senseID uuid.UUID) string {
	t.Helper()
	var s string
	if err := tx.QueryRow(t.Context(), `SELECT e.headword || ' ' || coalesce(e.noun_class, '-') || ' by ' || coalesce(i.author_id, '-')
		FROM senses s JOIN items i ON i.id = s.id JOIN entries e ON e.id = s.entry_id WHERE s.id = $1`, senseID).Scan(&s); err != nil {
		t.Fatal(err)
	}
	return s
}

// cardFor walks a reviewer's queue to the card for gloss.
func cardFor(ctx context.Context, t *testing.T, svc *review.Service, gloss string) (review.Card, bool) {
	t.Helper()
	var skip []uuid.UUID
	for {
		card, ok, err := svc.Next(ctx, "tst", skip)
		if err != nil {
			t.Fatal(err)
		}
		if !ok || card.GlossEn == gloss {
			return card, ok
		}
		skip = append(skip, card.ConceptID)
	}
}

func corrections(ids []uuid.UUID) []review.Decision {
	out := make([]review.Decision, len(ids))
	for i, id := range ids {
		out[i] = review.Decision{ItemID: id, Decision: correct}
	}
	return out
}

func TestCorrection(t *testing.T) {
	tx := dbtest.Tx(t)
	fixture(t, tx)
	svc := review.NewService(tx)
	stone := stoneCard(t, tx)

	got, err := svc.Decide(reviewer(t), corrections(stone), "", &review.Correction{Headword: "ilibwe", PartOfSpeech: "noun", NounClass: "5/6"})
	if err != nil {
		t.Fatal(err)
	}
	fix := got[0].CorrectionItemID
	if fix == nil || got[1].CorrectionItemID == nil || *got[1].CorrectionItemID != *fix {
		t.Fatalf("recorded = %+v, want both pointing at one new candidate", got)
	}
	if got := describe(t, tx, *fix); got != "ilibwe 5/6 by r1" {
		t.Errorf("new candidate = %s, want ilibwe 5/6 by r1", got)
	}

	// The corrector never sees their own word; another reviewer sees all three.
	if card, ok := cardFor(reviewer(t), t, svc, "zz stone"); ok {
		t.Errorf("r1 is offered zz stone again: %+v", card.Candidates)
	}
	if card, ok := cardFor(reviewer(t, "r2"), t, svc, "zz stone"); !ok || len(card.Candidates) != 3 {
		t.Errorf("r2 sees zz stone = %v with %d candidates, want 3", ok, len(card.Candidates))
	}

	// A second correction of the same candidate rejects it.
	if _, err := svc.Decide(reviewer(t, "r2"), corrections(stone[:1]), "", &review.Correction{Headword: "ilibwe"}); err != nil {
		t.Fatal(err)
	}
	item, _, err := items.NewStore(tx).Lock(t.Context(), stone[0])
	if err != nil {
		t.Fatal(err)
	}
	if item.Status != items.StatusRejected {
		t.Errorf("after two corrections: %s, want rejected", item.Status)
	}
}

func TestCorrectionRefused(t *testing.T) {
	tx := dbtest.Tx(t)
	fixture(t, tx)
	svc := review.NewService(tx)
	stone := stoneCard(t, tx)
	var candidate string
	if err := tx.QueryRow(t.Context(), `SELECT e.headword FROM senses s JOIN entries e ON e.id = s.entry_id WHERE s.id = $1`, stone[0]).Scan(&candidate); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name       string
		decisions  []review.Decision
		correction *review.Correction
	}{
		{"one of the candidates", corrections(stone), &review.Correction{Headword: candidate}},
		{"correct without a correction", corrections(stone), nil},
		{"a correction without correct", []review.Decision{{ItemID: stone[0], Decision: reject}}, &review.Correction{Headword: "ilibwe"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := svc.Decide(reviewer(t), tt.decisions, "", tt.correction); httpx.Status(err) != 400 {
				t.Errorf("err = %v, want 400", err)
			}
		})
	}
}
