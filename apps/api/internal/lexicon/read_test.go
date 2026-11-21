package lexicon_test

import (
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

// TestSearch uses made-up Bemba words, so the stemmer and style-guide
// spelling apply without colliding with real entries.
func TestSearch(t *testing.T) {
	lex, _ := searchFixture(t)
	ctx := t.Context()
	tests := []struct{ query, headword, match string }{
		{"UMUKÁZZO", "umukazzo", "headword"},
		{"ichipazzo", "icipazzo", "variant"},
		{"abakazzo", "umukazzo", "root"},
		{"umukazo", "umukazzo", "similar"},
		{"umupazzo", "umupazzo", ""}, // only a disputed sense: never returned
	}
	for _, tt := range tests {
		hits, err := lex.Search(ctx, "bem", tt.query, 50)
		if err != nil {
			t.Fatal(err)
		}
		h := find(hits, tt.headword)
		switch {
		case tt.match == "" && h != nil:
			t.Errorf("%q: %s returned, want left out", tt.query, tt.headword)
		case tt.match == "":
		case h == nil:
			t.Errorf("%q: %s not found", tt.query, tt.headword)
		case h.Match != tt.match || len(h.Entry.Senses) != 1 || h.Entry.Senses[0].Trust.Status != "machine" || len(h.Entry.Trust.Sources) != 1:
			t.Errorf("%q: hit = %s with %+v, want %s with one machine sense from Claude", tt.query, h.Match, h.Entry.Senses, tt.match)
		}
	}
}

func find(hits []lexicon.Hit, headword string) *lexicon.Hit {
	for i := range hits {
		if hits[i].Entry.Headword == headword {
			return &hits[i]
		}
	}
	return nil
}

func searchFixture(t *testing.T) (*lexicon.Store, pgx.Tx) {
	t.Helper()
	tx := dbtest.Tx(t)
	ctx := t.Context()
	claude, err := sources.Lookup(ctx, tx, "Claude")
	if err != nil {
		t.Fatal(err)
	}
	lex := lexicon.NewStore(tx)
	for i, w := range []struct{ headword, status string }{
		{"umukazzo", "machine"},
		{"ichipazzo", "machine"}, // stored as icipazzo, with ichipazzo as a variant
		{"umupazzo", "disputed"},
	} {
		added, err := lex.Add(ctx, lexicon.Proposal{
			Language: "bem", Headword: w.headword, GlossEn: "zz " + w.headword, Origin: items.Machine(claude, "m", string(rune('a'+i))),
		})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := tx.Exec(ctx, `UPDATE items SET status = $2 WHERE id = $1`, added.SenseID, w.status); err != nil {
			t.Fatal(err)
		}
	}

	return lex, tx
}

// TestGet reads a full entry: its stored variant, and its sense's example.
// An entry with only a disputed sense is not found.
func TestGet(t *testing.T) {
	lex, tx := searchFixture(t)
	ctx := t.Context()
	hits, err := lex.Search(ctx, "bem", "icipazzo", 50)
	if err != nil {
		t.Fatal(err)
	}
	hit := find(hits, "icipazzo")
	if hit == nil {
		t.Fatal("icipazzo not found")
	}
	if _, err := tx.Exec(ctx, `WITH i AS (INSERT INTO items (kind, language_code, status) VALUES ('example', 'bem', 'machine') RETURNING id)
		INSERT INTO examples (id, sense_id, text, translation_en) SELECT id, $1, 'zz text', 'zz translation' FROM i`, hit.Entry.Senses[0].ID); err != nil {
		t.Fatal(err)
	}
	e, ok, err := lex.Get(ctx, hit.Entry.ID)
	if err != nil || !ok {
		t.Fatalf("Get = %v, %v", ok, err)
	}
	if len(e.Variants) != 1 || e.Variants[0].Spelling != "ichipazzo" {
		t.Errorf("variants = %+v, want ichipazzo", e.Variants)
	}
	if s := e.Senses[0]; !s.ConceptID.Valid || len(s.Examples) != 1 || s.Examples[0].TranslationEn != "zz translation" || s.Examples[0].Trust.Status != "machine" {
		t.Errorf("sense = %+v, want a concept and the example", s)
	}

	var disputed uuid.UUID
	if err := tx.QueryRow(ctx, `SELECT id FROM entries WHERE language_code = 'bem' AND headword = 'umupazzo'`).Scan(&disputed); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := lex.Get(ctx, disputed); ok || err != nil {
		t.Errorf("disputed-only entry: Get = %v, %v; want not found", ok, err)
	}
}

// TestConceptEntries: a Nyanja word with the same meaning shares icipazzo's
// concept, so the concept lists both.
func TestConceptEntries(t *testing.T) {
	lex, tx := searchFixture(t)
	ctx := t.Context()
	claude, err := sources.Lookup(ctx, tx, "Claude")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := lex.Add(ctx, lexicon.Proposal{Language: "nya", Headword: "zzcinthu", GlossEn: "zz ichipazzo", Origin: items.Machine(claude, "m", "n")}); err != nil {
		t.Fatal(err)
	}
	hits, err := lex.Search(ctx, "bem", "icipazzo", 50)
	if err != nil {
		t.Fatal(err)
	}
	id := find(hits, "icipazzo").Entry.Senses[0].ConceptID.UUID
	c, ok, err := lex.ConceptEntries(ctx, id)
	if err != nil || !ok {
		t.Fatalf("ConceptEntries = %v, %v", ok, err)
	}
	var got []string
	for _, e := range c.Entries {
		got = append(got, e.Language+":"+e.Headword)
	}
	if c.GlossEn != "zz ichipazzo" || !slices.Equal(got, []string{"bem:icipazzo", "nya:zzcinthu"}) {
		t.Errorf("concept %q entries = %v", c.GlossEn, got)
	}
	if _, ok, err := lex.ConceptEntries(ctx, uuid.New()); ok || err != nil {
		t.Errorf("unknown concept: %v, %v; want not found", ok, err)
	}
}
