package httpx

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Optional returns nil for an empty string, so optional response fields are
// left out of the JSON instead of sent empty.
func Optional[T ~string](s T) *T {
	if s == "" {
		return nil
	}
	return &s
}

// Deref returns the value of an optional request field, or its zero value.
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// DerefOr returns the value of an optional request field, or def when absent.
func DerefOr[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}

// OptionalUUID returns a nullable ID as an optional response field.
func OptionalUUID(id uuid.NullUUID) *uuid.UUID {
	if !id.Valid {
		return nil
	}
	return &id.UUID
}

// OptionalTime returns nil for a NULL timestamp, so it is left out of JSON.
func OptionalTime(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
