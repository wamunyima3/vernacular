package languages

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"

	"github.com/wamunyima3/vernacular/apps/api/internal/languages/dbgen"
)

// Store reads languages from PostgreSQL.
type Store struct {
	q *dbgen.Queries
}

// NewStore wraps a pool or transaction.
func NewStore(db dbgen.DBTX) *Store {
	return &Store{q: dbgen.New(db)}
}

// List returns all languages in rollout order.
func (s *Store) List(ctx context.Context) ([]Language, error) {
	rows, err := s.q.ListLanguages(ctx)
	if err != nil {
		return nil, fmt.Errorf("list languages: %w", err)
	}
	out := make([]Language, len(rows))
	for i, r := range rows {
		out[i] = Language{Code: r.Code, Name: r.Name, OwnName: r.OwnName, Wave: int(r.Wave), Status: r.Status}
	}
	return out, nil
}

// RequireActive refuses a language that is unknown or planned but not yet
// open for data, as a bad request.
func (s *Store) RequireActive(ctx context.Context, code string) error {
	l, err := s.q.GetLanguageStatus(ctx, code)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return httpx.BadRequest(fmt.Sprintf("unknown language %q", code))
	case err != nil:
		return fmt.Errorf("get language %s: %w", code, err)
	case l.Status != string(gen.LanguageStatusActive):
		return httpx.BadRequest(l.Name + " is planned, not open yet")
	}
	return nil
}
