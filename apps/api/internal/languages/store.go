package languages

import (
	"context"
	"fmt"

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
