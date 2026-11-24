// Package items owns the shared `items` and `item_sources` tables
// (ADR-0006): every reviewable row is created here, with its trust status and
// where it came from. The review module (Phase 1) changes status through it.
package items

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

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

// Reviewable reports whether reviewers judge items of this kind: senses and
// phrases.
func Reviewable(kind dbgen.ItemKind) bool { return kind == KindSense || kind == KindPhrase }

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

// Contributed is the origin of rows a registered contributor proposed
// (CON-1). They start as `sourced`, like imports, and record their author.
func Contributed(sourceID uuid.UUID, authorID string) Origin {
	return Origin{Status: dbgen.TrustStatusSourced, SourceID: sourceID, Reference: "contribution", AuthorID: authorID}
}

// SourceRef is one source an item cites.
type SourceRef struct {
	Name, Licence, URL, Reference, Model string
}

// Trust is the provenance every returned item carries (PRV-1).
type Trust struct {
	Status               string
	Sources              []SourceRef
	Accepts, Corrections int
	Release              string // empty until the item is first released
}

// Trust reads an item's status, sources and review counts.
func (s *Store) Trust(ctx context.Context, id uuid.UUID) (Trust, error) {
	ts, err := s.Trusts(ctx, []uuid.UUID{id})
	if err != nil {
		return Trust{}, err
	}
	t, ok := ts[id]
	if !ok {
		return Trust{}, fmt.Errorf("get trust of %s: %w", id, pgx.ErrNoRows)
	}
	return t, nil
}

// Trusts reads the trust fields of many items in two queries; items that do
// not exist are left out.
func (s *Store) Trusts(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]Trust, error) {
	rows, err := s.q.ListTrusts(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("list trusts: %w", err)
	}
	refs, err := s.q.ListTrustSourcesOf(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("list trust sources: %w", err)
	}
	out := make(map[uuid.UUID]Trust, len(rows))
	for _, r := range rows {
		out[r.ID] = Trust{Status: string(r.Status), Accepts: int(r.Accepts), Corrections: int(r.Corrections), Release: r.Release}
	}
	for _, r := range refs {
		t := out[r.ItemID]
		t.Sources = append(t.Sources, SourceRef{
			Name: r.Name, Licence: r.Licence, URL: r.Url.String, Reference: r.Reference, Model: r.Model.String,
		})
		out[r.ItemID] = t
	}
	return out, nil
}

// Locked is an item held for a review decision.
type Locked struct {
	Kind     dbgen.ItemKind
	Language string
	Status   dbgen.TrustStatus
	AuthorID string
	Open     bool   // still takes reviews (open_for_review)
	Origin   Status // before any review: machine or sourced
}

// Lock reads an item and locks it until the transaction ends, so decisions
// on it are made one after another. False when there is no such item.
func (s *Store) Lock(ctx context.Context, id uuid.UUID) (Locked, bool, error) {
	row, err := s.q.LockItem(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Locked{}, false, nil
	}
	if err != nil {
		return Locked{}, false, fmt.Errorf("lock item %s: %w", id, err)
	}
	return Locked{Kind: row.Kind, Language: row.LanguageCode, Status: row.Status, AuthorID: row.AuthorID.String, Open: row.Open, Origin: row.Origin}, true, nil
}

// Status is an item's trust status (ADR-0004).
type Status = dbgen.TrustStatus

// Trust statuses. `rejected` is internal: never returned by the API.
const (
	StatusVerified = dbgen.TrustStatusVerified
	StatusReviewed = dbgen.TrustStatusReviewed
	StatusSourced  = dbgen.TrustStatusSourced
	StatusMachine  = dbgen.TrustStatusMachine
	StatusDisputed = dbgen.TrustStatusDisputed
	StatusRejected = dbgen.TrustStatusRejected
)

// SetStatus records a status the review module's promotion function decided.
// Nothing else may call it (AGENTS.md rule 2).
func (s *Store) SetStatus(ctx context.Context, id uuid.UUID, status Status) error {
	if err := s.q.SetItemStatus(ctx, dbgen.SetItemStatusParams{ID: id, Status: status}); err != nil {
		return fmt.Errorf("set status of %s: %w", id, err)
	}
	return nil
}
