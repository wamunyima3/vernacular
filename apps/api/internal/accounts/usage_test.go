package accounts_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// newKey adds user u1 with one API key.
func newKey(t *testing.T, tx pgx.Tx) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(t.Context(), `
		WITH u AS (INSERT INTO users (id, name, email, email_verified, role) VALUES ('u1', 'U', 'u1@test.invalid', true, 'developer') RETURNING id)
		INSERT INTO api_keys (user_id, name, prefix, hash) SELECT id, 'app', 'vk_x', 'h' FROM u RETURNING id`).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

// TestKeyUsage: calls with a key are counted by outcome, flushed, and shown
// only to the key's owner (DEV-5, US-13).
func TestKeyUsage(t *testing.T) {
	tx := dbtest.Tx(t)
	key := newKey(t, tx)
	usage := accounts.NewUsage(tx)
	for _, status := range []int{200, 200, 404, 429} {
		h := usage.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(status) }))
		ctx := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "u1", KeyID: key})
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/languages", nil)
		h.ServeHTTP(httptest.NewRecorder(), req)
	}
	for range 2 { // a second flush with nothing new changes nothing
		if err := usage.Flush(t.Context()); err != nil {
			t.Fatal(err)
		}
	}

	svc := accounts.NewService(tx)
	owner := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "u1", Role: gen.RoleDeveloper})
	days, err := svc.KeyUsage(owner, key, 30)
	if err != nil || len(days) != 1 {
		t.Fatalf("usage = %+v, %v; want one day", days, err)
	}
	if d := days[0]; d.Calls != 4 || d.Errors != 1 || d.Limited != 1 {
		t.Errorf("today = %+v, want 4 calls, 1 error, 1 limited", d)
	}
	other := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "u2", Role: gen.RoleDeveloper})
	if _, err := svc.KeyUsage(other, key, 30); httpx.Status(err) != 404 {
		t.Errorf("another user's key: %v, want 404", err)
	}
}
