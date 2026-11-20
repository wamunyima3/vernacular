package accounts

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts/dbgen"
)

// ErrNoSession means the token is unknown, expired or belongs to a banned user.
var ErrNoSession = errors.New("no valid session")

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
	rows, err := s.q.ListApprovals(ctx, u.ID)
	if err != nil {
		return Principal{}, fmt.Errorf("list approvals: %w", err)
	}
	p := Principal{ID: u.ID, Name: u.Name, Email: u.Email, Role: u.Role, Approvals: make([]Approval, len(rows))}
	for i, r := range rows {
		p.Approvals[i] = Approval{Language: r.LanguageCode, Level: string(r.Level)}
	}
	return p, nil
}
