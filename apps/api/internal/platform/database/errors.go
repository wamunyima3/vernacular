package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// SQLSTATE codes callers turn into API errors.
const (
	UniqueViolation     = "23505"
	ForeignKeyViolation = "23503"
)

// Violates reports whether err is a Postgres error with the given SQLSTATE.
func Violates(err error, code string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == code
}
