package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Conn is a pool or a transaction: what stores query through.
type Conn interface {
	Beginner
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type txKey struct{}

// WithTx makes every Routed handle use tx for the rest of the request. A
// pinned read uses it to run the live queries on a release's frozen rows
// (ADR-0012).
func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// Routed runs each query on the request's transaction when WithTx set one,
// and on Default otherwise. Read stores are built over it once, at startup.
type Routed struct{ Default Conn }

func (r Routed) conn(ctx context.Context) Conn {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return r.Default
}

// Exec runs sql on the request's connection.
func (r Routed) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return r.conn(ctx).Exec(ctx, sql, args...)
}

// Query runs sql on the request's connection.
func (r Routed) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return r.conn(ctx).Query(ctx, sql, args...)
}

// QueryRow runs sql on the request's connection.
func (r Routed) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return r.conn(ctx).QueryRow(ctx, sql, args...)
}

// Begin starts a transaction, or a savepoint inside the request's one.
func (r Routed) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.conn(ctx).Begin(ctx)
}
