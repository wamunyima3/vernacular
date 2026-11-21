package accounts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// ErrNoSession means the token is unknown, expired or belongs to a banned
// user. It wraps httpx.ErrUnauthorized, so it answers 401.
var ErrNoSession = fmt.Errorf("no valid session: %w", httpx.ErrUnauthorized)

// Store reads Better Auth's session and user tables.
type Store struct {
	q *dbgen.Queries
}

// NewStore wraps a pool or transaction.
func NewStore(db dbgen.DBTX) *Store {
	return &Store{q: dbgen.New(db)}
}

// Authenticate returns the user behind a session token.
func (s *Store) Authenticate(ctx context.Context, token string) (Principal, error) {
	u, err := s.q.GetSessionUser(ctx, token)
	if errors.Is(err, pgx.ErrNoRows) {
		return Principal{}, ErrNoSession
	}
	if err != nil {
		return Principal{}, fmt.Errorf("get session user: %w", err)
	}
	p := Principal{ID: u.ID, Name: u.Name, Email: u.Email, Role: gen.Role(u.Role), Credited: u.Credited}
	if err := json.Unmarshal(u.Approvals, &p.Approvals); err != nil {
		return Principal{}, fmt.Errorf("decode approvals: %w", err)
	}
	return p, nil
}
