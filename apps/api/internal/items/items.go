// Package items owns the shared `items` and `item_sources` tables
// (ADR-0006): every reviewable row is created here, with its trust status and
// where it came from. The review module (Phase 1) changes status through it.
package items

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/items/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
)

// Kinds of reviewable rows.
const (
	KindEntry   = dbgen.ItemKindEntry
	KindSense   = dbgen.ItemKindSense
	KindExample = dbgen.ItemKindExample
	KindPhrase  = dbgen.ItemKindPhrase
)

// Origin says where a new item came from and how far it is trusted.
type Origin struct {
	Status    dbgen.TrustStatus
	SourceID  uuid.UUID
	Reference string // where in the source
	Model     string // model version, for machine items (AI-2)
	AuthorID  string // contributor, when a person proposed it
}

// Sourced is the origin of rows imported from a licensed source.
func Sourced(sourceID uuid.UUID, reference string) Origin {
	return Origin{Status: dbgen.TrustStatusSourced, SourceID: sourceID, Reference: reference}
}

// Machine is the origin of candidates a model proposed (AI-1, AI-2). The
// source row names the provider; model is the exact model version.
func Machine(sourceID uuid.UUID, model, reference string) Origin {
	return Origin{Status: dbgen.TrustStatusMachine, SourceID: sourceID, Reference: reference, Model: model}
}

// Store writes items. Give it the transaction the typed row is written in.
type Store struct {
	q *dbgen.Queries
}

// NewStore wraps a pool or transaction.
func NewStore(db dbgen.DBTX) *Store { return &Store{q: dbgen.New(db)} }

// Create makes the items row, runs create to insert the typed row with the
// same ID, and links the source. It is the only way reviewable rows are made.
func (s *Store) Create(ctx context.Context, kind dbgen.ItemKind, language string, o Origin, create func(uuid.UUID) error) (uuid.UUID, error) {
	id, err := s.q.CreateItem(ctx, dbgen.CreateItemParams{
		Kind: kind, LanguageCode: language, Status: o.Status, AuthorID: database.Text(o.AuthorID),
	})
	if err != nil {
		return id, fmt.Errorf("create item: %w", err)
	}
	if err := create(id); err != nil {
		return id, err
	}
	return id, s.Link(ctx, id, o)
}

// Link adds an origin to an existing item's sources without changing its
// status. Linking the same origin twice changes nothing.
func (s *Store) Link(ctx context.Context, id uuid.UUID, o Origin) error {
	err := s.q.LinkSource(ctx, dbgen.LinkSourceParams{ItemID: id, SourceID: o.SourceID, Reference: o.Reference, Model: database.Text(o.Model)})
	if err != nil {
		return fmt.Errorf("link source: %w", err)
	}
	return nil
}
