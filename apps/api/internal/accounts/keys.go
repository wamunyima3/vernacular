package accounts

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// Developer API keys (DEV-1, NFR-5, US-01). A key is "vk_" and 52 random
// base32 characters (256 bits); only its SHA-256 is stored. A fast hash is
// enough because, unlike a password, the key cannot be guessed.

// ErrNoKey means the API key is unknown, revoked or belongs to a banned user.
var ErrNoKey = fmt.Errorf("unknown or revoked API key: %w", httpx.ErrUnauthorized)

const (
	keyPrefix  = "vk_"
	shownChars = len(keyPrefix) + 4 // how much of a key the owner sees again
)

func hashKey(secret string) []byte {
	h := sha256.Sum256([]byte(secret))
	return h[:]
}

// AuthenticateKey returns the owner of a live API key.
func (s *Store) AuthenticateKey(ctx context.Context, secret string) (Principal, error) {
	if !strings.HasPrefix(secret, keyPrefix) {
		return Principal{}, ErrNoKey
	}
	k, err := s.q.GetKeyOwner(ctx, hashKey(secret))
	if errors.Is(err, pgx.ErrNoRows) {
		return Principal{}, ErrNoKey
	}
	if err != nil {
		return Principal{}, fmt.Errorf("get key owner: %w", err)
	}
	if !k.LastUsedAt.Valid || time.Since(k.LastUsedAt.Time) > time.Hour {
		if err := s.q.TouchKey(ctx, k.ID); err != nil {
			return Principal{}, fmt.Errorf("touch key: %w", err)
		}
	}
	return Principal{ID: k.UserID, KeyID: k.ID}, nil
}

// Keys lists the caller's API keys, live first.
func (s *Service) Keys(ctx context.Context) ([]dbgen.ListKeysRow, error) {
	me, err := Require(ctx, nil)
	if err != nil {
		return nil, err
	}
	var rows []dbgen.ListKeysRow
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		rows, err = dbgen.New(tx).ListKeys(ctx, me.ID)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("list keys: %w", err)
	}
	return rows, nil
}

// NewKey is a created key: its secret is returned only this once.
type NewKey struct {
	dbgen.ListKeysRow
	Secret string
}

// CreateKey makes an API key for the caller. At most 10 are live at once.
func (s *Service) CreateKey(ctx context.Context, name string) (NewKey, error) {
	me, err := Require(ctx, nil)
	if err != nil {
		return NewKey{}, err
	}
	if name = strings.TrimSpace(name); name == "" {
		return NewKey{}, httpx.BadRequest("Give the key a name")
	}
	secret := keyPrefix + rand.Text() + rand.Text()
	k := NewKey{Secret: secret, ListKeysRow: dbgen.ListKeysRow{Name: name, Prefix: secret[:shownChars]}}
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		row, err := dbgen.New(tx).CreateKey(ctx, dbgen.CreateKeyParams{UserID: me.ID, Name: name, Prefix: k.Prefix, Hash: hashKey(secret)})
		if errors.Is(err, pgx.ErrNoRows) {
			return httpx.Conflict("You have 10 live keys. Revoke one first.")
		}
		k.ID, k.CreatedAt = row.ID, row.CreatedAt
		return err
	})
	if err != nil {
		return NewKey{}, fmt.Errorf("create key: %w", err)
	}
	return k, nil
}

// RevokeKey stops one of the caller's keys working, from now on.
func (s *Service) RevokeKey(ctx context.Context, id uuid.UUID) error {
	me, err := Require(ctx, nil)
	if err != nil {
		return err
	}
	return database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		n, err := dbgen.New(tx).RevokeKey(ctx, dbgen.RevokeKeyParams{ID: id, UserID: me.ID})
		if err != nil {
			return fmt.Errorf("revoke key: %w", err)
		}
		if n == 0 {
			return httpx.ErrNotFound
		}
		return nil
	})
}
