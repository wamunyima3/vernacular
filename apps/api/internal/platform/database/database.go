// Package database opens the shared PostgreSQL connection pool.
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect opens a pool and checks that the database answers.
func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("open pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}

// Beginner starts a transaction. A pool works, and so does an outer
// transaction, which makes the inner one a savepoint (used by tests).
type Beginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// InTx runs fn in a transaction with the audit actor set (NFR-8), committing
// only if fn succeeds. Every write path goes through it, so every audit row
// names who made the change.
func InTx(ctx context.Context, db Beginner, actor string, fn func(pgx.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, "SELECT set_config('app.actor_id', $1, true)", actor); err != nil {
		return fmt.Errorf("set audit actor: %w", err)
	}
	if err := fn(tx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}
