// Integration tests for rules the schema itself enforces (P0-ENG-01).
package db_test

import (
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
)

func TestFold(t *testing.T) {
	tx := dbtest.Tx(t)
	cases := map[string]string{
		"Amenshi":   "amenshi",
		"Ichíbemba": "ichibemba",
		"ŋombe":     "ng'ombe",
		"NG’OMBE":   "ng'ombe",
		"aŵiri":     "awiri",
		"cinthu":    "cinthu",
	}
	for in, want := range cases {
		var got string
		if err := tx.QueryRow(t.Context(), "SELECT fold($1)", in).Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("fold(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAuditRecordsActorAndChange(t *testing.T) {
	tx := dbtest.Tx(t)
	exec(t, tx, "SELECT set_config('app.actor_id', 'user-1', true)")
	var id string
	if err := tx.QueryRow(t.Context(),
		"INSERT INTO concepts (gloss_en) VALUES ('water') RETURNING id").Scan(&id); err != nil {
		t.Fatal(err)
	}
	exec(t, tx, "UPDATE concepts SET domain = 'nature' WHERE id = $1", id)

	var actor, before, after string
	err := tx.QueryRow(t.Context(), `
		SELECT actor_id, before->>'domain', after->>'domain' FROM audit_log
		WHERE table_name = 'concepts' AND row_id = $1 AND action = 'update'`, id).
		Scan(&actor, &before, &after)
	if err != nil {
		t.Fatal(err)
	}
	if actor != "user-1" || before != "general" || after != "nature" {
		t.Errorf("audit = (%s, %s, %s), want (user-1, general, nature)", actor, before, after)
	}
}

func TestConstraintsRejectBadRows(t *testing.T) {
	cases := map[string]string{
		"decomposed headword": `
			WITH i AS (INSERT INTO items (kind, language_code, status) VALUES ('entry', 'bem', 'sourced') RETURNING id)
			INSERT INTO entries (id, language_code, headword, part_of_speech)
			SELECT id, 'bem', U&'ici\0301ntu', 'noun' FROM i`, // i + combining acute, not NFC
		"unknown variety": `
			WITH i AS (INSERT INTO items (kind, language_code, status) VALUES ('entry', 'bem', 'sourced') RETURNING id)
			INSERT INTO entries (id, language_code, headword, part_of_speech, variety)
			SELECT id, 'bem', 'umuntu', 'noun', 'town' FROM i`,
		"correction without changes": `
			WITH i AS (INSERT INTO items (kind, language_code, status) VALUES ('entry', 'bem', 'sourced') RETURNING id)
			INSERT INTO reviews (item_id, reviewer_id, decision) SELECT id, 'r1', 'correct' FROM i`,
		"bad release name": `INSERT INTO releases (name) VALUES ('v1')`,
	}
	for name, sql := range cases {
		t.Run(name, func(t *testing.T) {
			tx := dbtest.Tx(t)
			if _, err := tx.Exec(t.Context(), sql); err == nil {
				t.Error("insert succeeded, want a constraint error")
			}
		})
	}
}

func exec(t *testing.T, tx pgx.Tx, sql string, args ...any) {
	t.Helper()
	if _, err := tx.Exec(t.Context(), sql, args...); err != nil {
		t.Fatal(err)
	}
}
