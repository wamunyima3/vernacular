package sources_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

// Made-up words, so the test does not depend on what `make seed` imported.
const sample = `{"src": "test water", "trgs": ["zzamenshi"]}
{"src": "Test hello", "trgs": ["Zzmuli shani", "zzhalo"]}
{"src": "test hi", "trgs": ["Zzmuli shani"]}
{"src": "test thing", "trgs": ["zzichintu", "zzicintu"]}`

func feed(t *testing.T) sources.Feed {
	t.Helper()
	f, ok := sources.FeedByName("gatitos-en-bem")
	if !ok {
		t.Fatal("feed gatitos-en-bem missing")
	}
	return f
}

func TestIngestIsIdempotent(t *testing.T) {
	tx := dbtest.Tx(t)
	f := feed(t)

	first, err := sources.Ingest(t.Context(), tx, f, strings.NewReader(sample))
	if err != nil {
		t.Fatal(err)
	}
	// "zzmuli shani" is one entry with two senses (hello, hi). "zzichintu" is
	// stored as the standard "zzicintu" with a variant, and the source's own
	// "zzicintu" then merges into that entry.
	if want := (sources.Stats{Records: 6, Entries: 4, Senses: 5, Variants: 1}); first != want {
		t.Errorf("first run = %+v, want %+v", first, want)
	}
	again, err := sources.Ingest(t.Context(), tx, f, strings.NewReader(sample))
	if err != nil {
		t.Fatal(err)
	}
	if want := (sources.Stats{Records: 6}); again != want {
		t.Errorf("second run = %+v, want %+v (nothing new)", again, want)
	}

	var status, actor, licence string
	err = tx.QueryRow(t.Context(), `
		SELECT i.status, a.actor_id, s.licence
		FROM entries e
		JOIN items i ON i.id = e.id
		JOIN item_sources isrc ON isrc.item_id = e.id
		JOIN sources s ON s.id = isrc.source_id
		JOIN audit_log a ON a.table_name = 'entries' AND a.row_id = e.id::text
		WHERE e.headword = 'zzamenshi'`).Scan(&status, &actor, &licence)
	if err != nil {
		t.Fatal(err)
	}
	if status != "sourced" || actor != "ingest:gatitos-en-bem" || licence != "CC BY 4.0" {
		t.Errorf("zzamenshi = (%s, %s, %s)", status, actor, licence)
	}
}

func TestIngestRefusesUnregisteredSource(t *testing.T) {
	tx := dbtest.Tx(t)
	f := feed(t)
	f.Source = "Unlicensed Dictionary"
	_, err := sources.Ingest(t.Context(), tx, f, strings.NewReader(sample))
	if !errors.Is(err, sources.ErrNotRegistered) {
		t.Fatalf("err = %v, want ErrNotRegistered", err)
	}
}
