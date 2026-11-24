package accounts_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// TestKeys: u1 creates a key, calls with it, then revokes it (US-01).
func TestKeys(t *testing.T) {
	tx := dbtest.Tx(t)
	agreementFixture(t, tx)
	svc, store := accounts.NewService(tx), accounts.NewStore(tx)
	me := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "u1"})

	k := createKey(me, t, tx, svc)

	p, err := store.AuthenticateKey(t.Context(), k.Secret)
	if err != nil || p.ID != "u1" || p.KeyID != k.ID {
		t.Fatalf("authenticate = %+v, %v", p, err)
	}
	keys, err := svc.Keys(me)
	if err != nil || len(keys) != 1 || !keys[0].LastUsedAt.Valid {
		t.Fatalf("keys = %+v, %v; want one, used", keys, err)
	}

	if err := svc.RevokeKey(me, k.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := store.AuthenticateKey(t.Context(), k.Secret); !errors.Is(err, accounts.ErrNoKey) {
		t.Errorf("revoked key: %v, want ErrNoKey", err)
	}
	refused := []struct {
		name string
		err  error
		want int
	}{
		{"a key lists keys", func() error { _, err := svc.Keys(accounts.WithPrincipal(t.Context(), p)); return err }(), 401},
		{"revoking twice", svc.RevokeKey(me, k.ID), 404},
		{"someone else's key", svc.RevokeKey(accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "u2"}), k.ID), 404},
	}
	for _, tt := range refused {
		if httpx.Status(tt.err) != tt.want {
			t.Errorf("%s: %v, want %d", tt.name, tt.err, tt.want)
		}
	}
}

func TestKeyLimit(t *testing.T) {
	tx := dbtest.Tx(t)
	agreementFixture(t, tx)
	svc := accounts.NewService(tx)
	me := accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "u1"})
	for i := range 10 {
		if _, err := svc.CreateKey(me, fmt.Sprint("key ", i)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.CreateKey(me, "eleventh"); httpx.Status(err) != 409 {
		t.Errorf("an 11th live key: %v, want 409", err)
	}
	if _, err := svc.CreateKey(me, "  "); httpx.Status(err) != 400 {
		t.Errorf("no name: %v, want 400", err)
	}
}

// createKey makes u1's key and checks that only its hash is stored.
func createKey(me context.Context, t *testing.T, tx pgx.Tx, svc *accounts.Service) accounts.NewKey {
	t.Helper()
	k, err := svc.CreateKey(me, " My SMS app ")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(k.Secret, k.Prefix) || len(k.Secret) != 55 || k.Name != "My SMS app" {
		t.Fatalf("key = %+v", k)
	}
	var stored int
	if err := tx.QueryRow(t.Context(), `SELECT count(*) FROM api_keys WHERE hash = sha256($1::bytea)`, k.Secret).Scan(&stored); err != nil || stored != 1 {
		t.Fatalf("stored by hash: %d, %v", stored, err)
	}
	return k
}
