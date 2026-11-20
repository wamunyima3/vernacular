package database

import "github.com/jackc/pgx/v5/pgtype"

// Text stores an empty string as NULL.
func Text(s string) pgtype.Text { return pgtype.Text{String: s, Valid: s != ""} }
