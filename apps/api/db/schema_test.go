// Integration tests for rules the schema itself enforces (P0-ENG-01).
package db_test

import (
	"slices"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
)

func TestAuditRecordsActorAndChange(t *testing.T) {
	tx := dbtest.Tx(t)
	exec(t, tx, "SELECT set_config('app.actor_id', 'user-1', true)")
	var id string
	if err := tx.QueryRow(t.Context(),
		"INSERT INTO concepts (gloss_en) VALUES ('zz audit test') RETURNING id").Scan(&id); err != nil {
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

// The contract's TrustStatus and the database enum must list the same values.
func TestTrustStatusMatchesContract(t *testing.T) {
	tx := dbtest.Tx(t)
	var db []string
	if err := tx.QueryRow(t.Context(), "SELECT enum_range(NULL::trust_status)::text[]").Scan(&db); err != nil {
		t.Fatal(err)
	}
	contract := []string{
		string(gen.TrustStatusVerified), string(gen.TrustStatusReviewed), string(gen.TrustStatusSourced),
		string(gen.TrustStatusMachine), string(gen.TrustStatusDisputed),
		"rejected", // internal, never returned (ADR-0004 amendment)
	}
	if !slices.Equal(db, contract) {
		t.Errorf("database trust_status = %v, contract = %v", db, contract)
	}
}

func exec(t *testing.T, tx pgx.Tx, sql string, args ...any) {
	t.Helper()
	if _, err := tx.Exec(t.Context(), sql, args...); err != nil {
		t.Fatal(err)
	}
}
