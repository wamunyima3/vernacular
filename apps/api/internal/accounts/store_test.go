package accounts_test

import (
	"errors"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
)

// Rows are written in the shape Better Auth writes them (ADR-0007), so a
// Better Auth upgrade that changes the shape fails here.
func TestStoreAuthenticate(t *testing.T) {
	tx := dbtest.Tx(t)
	for _, sql := range []string{
		`INSERT INTO users (id, name, email, email_verified, role) VALUES
			('u1', 'Mwila', 'mwila@test.invalid', true, 'contributor'),
			('u2', 'Banned', 'banned@test.invalid', true, DEFAULT)`,
		`UPDATE users SET banned = true WHERE id = 'u2'`,
		`INSERT INTO sessions (id, token, user_id, expires_at, updated_at) VALUES
			('s1', 'live', 'u1', now() + interval '1 day', now()),
			('s2', 'expired', 'u1', now() - interval '1 second', now()),
			('s3', 'banned', 'u2', now() + interval '1 day', now())`,
		`INSERT INTO reviewer_approvals (user_id, language_code, level) VALUES ('u1', 'bem', 'senior')`,
	} {
		if _, err := tx.Exec(t.Context(), sql); err != nil {
			t.Fatal(err)
		}
	}
	store := accounts.NewStore(tx)

	p, err := store.Authenticate(t.Context(), "live")
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != "u1" || p.Role != gen.RoleContributor || !p.IsSenior("bem") {
		t.Errorf("principal = %+v", p)
	}

	for _, token := range []string{"expired", "banned", "unknown"} {
		if _, err := store.Authenticate(t.Context(), token); !errors.Is(err, accounts.ErrNoSession) {
			t.Errorf("%s: err = %v, want ErrNoSession", token, err)
		}
	}
}
