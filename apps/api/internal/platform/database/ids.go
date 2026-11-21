package database

import "github.com/google/uuid"

// IDs collects the IDs of rows, for a follow-up `= ANY(@ids)` query.
func IDs[T any](rows []T, id func(T) uuid.UUID) []uuid.UUID {
	out := make([]uuid.UUID, len(rows))
	for i, r := range rows {
		out[i] = id(r)
	}
	return out
}
