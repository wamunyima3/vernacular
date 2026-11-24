// Package releases freezes the verified items into numbered dataset releases
// with counts and a change log (PRV-4, US-29).
package releases

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/releases/dbgen"
)

// Release is one frozen copy of the verified data.
type Release struct {
	Name        string
	PublishedAt time.Time
	Notes       string
	// Languages counts items per language code, then per kind (entry, sense, example, phrase).
	Languages map[string]map[string]int
	Changes   Changes
}

// Changes compares a release with the one before it.
type Changes struct {
	Added   int `json:"added"`
	Changed int `json:"changed"`
	Removed int `json:"removed"`
}

// counts is the shape SummariseRelease stores in releases.counts.
type counts struct {
	Languages map[string]map[string]int `json:"languages"`
	Changes
}

// DB is what the service needs from the pool.
type DB interface {
	database.Beginner
	dbgen.DBTX
}

// Service lists and builds releases.
type Service struct {
	db DB
	q  *dbgen.Queries
}

// NewService returns a release service over the pool.
func NewService(db DB) *Service {
	return &Service{db: db, q: dbgen.New(db)}
}

// List returns every release, newest first.
func (s *Service) List(ctx context.Context) ([]Release, error) {
	rows, err := s.q.ListReleases(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Release, len(rows))
	for i, r := range rows {
		if out[i], err = fromRow(r); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// Build freezes every verified item into a new release (admin only). It is
// refused when nothing is verified, or nothing changed since the last release.
func (s *Service) Build(ctx context.Context, notes string) (Release, error) {
	actor, err := accounts.RequireAdmin(ctx)
	if err != nil {
		return Release{}, err
	}
	var built Release
	err = database.InTx(ctx, s.db, actor.ID, func(tx pgx.Tx) error {
		q := dbgen.New(tx)
		if err := q.LockReleases(ctx); err != nil {
			return err
		}
		all, err := q.ListReleases(ctx)
		if err != nil {
			return err
		}
		var prev dbgen.Release // zero when this is the first release
		if len(all) > 0 {
			prev = all[0]
		}
		row, err := q.CreateRelease(ctx, notes)
		if err != nil {
			return err
		}
		n, err := q.FreezeVerified(ctx, row.ID)
		if err != nil {
			return err
		}
		if n == 0 {
			return httpx.Conflict("Nothing is verified yet, so there is nothing to release")
		}
		row, err = q.SummariseRelease(ctx, dbgen.SummariseReleaseParams{
			ReleaseID:  row.ID,
			PreviousID: uuid.NullUUID{UUID: prev.ID, Valid: prev.ID != uuid.Nil},
		})
		if err != nil {
			return err
		}
		if built, err = fromRow(row); err != nil {
			return err
		}
		if prev.ID != uuid.Nil && built.Changes == (Changes{}) {
			return httpx.Conflict(fmt.Sprintf("Nothing changed since release %s", prev.Name))
		}
		// The copy pinned reads use (ADR-0012).
		return q.MaterialiseRelease(ctx, dbgen.MaterialiseReleaseParams{ReleaseID: row.ID, Schema: Schema(row.Name)})
	})
	return built, err
}

// Schema is the PostgreSQL schema a release is materialised in, such as
// release_2027_1 for 2027.1 (ADR-0012).
func Schema(name string) string { return "release_" + strings.ReplaceAll(name, ".", "_") }

func fromRow(r dbgen.Release) (Release, error) {
	var c counts
	if err := json.Unmarshal(r.Counts, &c); err != nil {
		return Release{}, fmt.Errorf("release %s counts: %w", r.Name, err)
	}
	return Release{
		Name: r.Name, PublishedAt: r.PublishedAt.Time, Notes: r.Changelog,
		Languages: c.Languages, Changes: c.Changes,
	}, nil
}
