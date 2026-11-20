// Package dbtest gives integration tests a real database transaction that is
// rolled back when the test ends, so tests never see each other's rows.
package dbtest

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
)

// Tx opens a transaction on DATABASE_URL (migrated with `make migrate`).
// The test is skipped when DATABASE_URL is unset; `make test-api` and CI set it.
func Tx(t *testing.T) pgx.Tx {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set; run through `make test-api`")
	}
	pool, err := database.Connect(t.Context(), url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	tx, err := pool.Begin(t.Context())
	if err != nil {
		pool.Close()
		t.Fatalf("begin: %v", err)
	}
	// t.Context() is cancelled before cleanup runs.
	ctx := context.WithoutCancel(t.Context())
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		pool.Close()
	})
	return tx
}
