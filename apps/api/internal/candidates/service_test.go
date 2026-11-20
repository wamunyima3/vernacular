package candidates_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/candidates"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

type fixed candidates.Result

func (f fixed) Generate(context.Context, candidates.Request) (candidates.Result, error) {
	return candidates.Result(f), nil
}

func TestGenerateStoresMachineCandidates(t *testing.T) {
	ctx := context.Background()
	tx := dbtest.Tx(t)

	// A concept a licensed source already gives one word for.
	lex := lexicon.NewStore(tx)
	gatitos, err := sources.Lookup(ctx, tx, "GATITOS")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := lex.Add(ctx, lexicon.Proposal{
		Language: "bem", Headword: "zzamenshi", GlossEn: "zz water", Core: true, Origin: items.Sourced(gatitos, "test"),
	}); err != nil {
		t.Fatal(err)
	}
	var conceptID uuid.UUID
	if err := tx.QueryRow(ctx, "SELECT id FROM concepts WHERE gloss_en = 'zz water'").Scan(&conceptID); err != nil {
		t.Fatal(err)
	}

	claude, err := sources.Lookup(ctx, tx, "Claude")
	if err != nil {
		t.Fatal(err)
	}
	gen := fixed{Model: "claude-test-1", Candidates: []candidates.Candidate{
		{Headword: "zzamenshi", PartOfSpeech: ""},                                                 // agrees with the source
		{Headword: "zzichimenshi", PartOfSpeech: "noun", NounClass: "7/8", Plural: "zzifimenshi"}, // new, non-standard spelling
		{Headword: "zzbad", PartOfSpeech: "noun", NounClass: "seven"},                             // malformed class is dropped
	}}
	svc := candidates.NewService(tx, gen, claude, "ai:test")

	stats, err := svc.Generate(ctx, conceptID, "bem")
	if err != nil {
		t.Fatal(err)
	}
	if want := (candidates.Stats{Proposed: 3, Senses: 2}); stats != want {
		t.Errorf("first run = %+v, want %+v", stats, want)
	}
	again, err := svc.Generate(ctx, conceptID, "bem")
	if err != nil {
		t.Fatal(err)
	}
	if want := (candidates.Stats{Proposed: 3}); again != want {
		t.Errorf("second run = %+v, want %+v", again, want)
	}

	rows, err := tx.Query(ctx, `
		SELECT e.headword, coalesce(e.noun_class, ''), coalesce(e.plural, ''), i.status::text, (SELECT ra.priority FROM review_agreement ra WHERE ra.item_id = s.id),
		       (SELECT string_agg(coalesce(isrc.model, '-'), ',' ORDER BY isrc.model) FROM item_sources isrc WHERE isrc.item_id = s.id)
		FROM senses s JOIN entries e ON e.id = s.entry_id JOIN items i ON i.id = s.id
		WHERE s.concept_id = $1 ORDER BY e.headword`, conceptID)
	if err != nil {
		t.Fatal(err)
	}
	type row struct{ headword, class, plural, status, priority, models string }
	got, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (row, error) {
		var x row
		return x, r.Scan(&x.headword, &x.class, &x.plural, &x.status, &x.priority, &x.models)
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []row{
		// Status kept, Claude added as evidence; both sources agree, so it is a quick confirm.
		{"zzamenshi", "", "", "sourced", "quick", "claude-test-1,-"},
		{"zzbad", "", "", "machine", "full", "claude-test-1"},
		{"zzicimenshi", "7/8", "zzifimenshi", "machine", "full", "claude-test-1"},
	}
	if len(got) != len(want) {
		t.Fatalf("senses = %+v, want %+v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("sense %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// A second provider changes which answers agree: now three sources have
// answered, so a word two of them give is a standard review.
func TestAgreementAcrossProviders(t *testing.T) {
	ctx := context.Background()
	tx := dbtest.Tx(t)
	lex := lexicon.NewStore(tx)
	if _, err := lex.AddCoreConcept(ctx, "zz fire", "everyday"); err != nil {
		t.Fatal(err)
	}
	var conceptID uuid.UUID
	if err := tx.QueryRow(ctx, "SELECT id FROM concepts WHERE gloss_en = 'zz fire'").Scan(&conceptID); err != nil {
		t.Fatal(err)
	}
	run := func(source string, words ...string) {
		t.Helper()
		id, err := sources.Lookup(ctx, tx, source)
		if err != nil {
			t.Fatal(err)
		}
		gen := fixed{Model: source + "-test"}
		for _, w := range words {
			gen.Candidates = append(gen.Candidates, candidates.Candidate{Headword: w})
		}
		if _, err := candidates.NewService(tx, gen, id, "ai:test").Generate(ctx, conceptID, "bem"); err != nil {
			t.Fatal(err)
		}
	}
	priorities := func() map[string]string {
		t.Helper()
		rows, err := tx.Query(ctx, `SELECT e.headword, ra.priority FROM review_agreement ra
			JOIN senses s ON s.id = ra.item_id JOIN entries e ON e.id = s.entry_id WHERE ra.concept_id = $1`, conceptID)
		if err != nil {
			t.Fatal(err)
		}
		got := map[string]string{}
		for rows.Next() {
			var w, p string
			if err := rows.Scan(&w, &p); err != nil {
				t.Fatal(err)
			}
			got[w] = p
		}
		return got
	}

	run("Claude", "zzmulilo")
	if got := priorities(); got["zzmulilo"] != "full" {
		t.Errorf("one source alone = %v, want full", got)
	}
	run("Google Translate", "zzmulilo")
	if got := priorities(); got["zzmulilo"] != "quick" {
		t.Errorf("both providers agree = %v, want quick", got)
	}
	run("GATITOS", "zzmoto")
	if got := priorities(); got["zzmulilo"] != "standard" || got["zzmoto"] != "full" {
		t.Errorf("two of three agree = %v, want zzmulilo standard, zzmoto full", got)
	}
}
