// Package sources owns the source register (DQ-1) and ingestion: one parser
// per licensed source, all feeding the lexicon through the same path.
package sources

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources/dbgen"
)

// ErrNotRegistered means the source has no row in the register, so nothing
// from it may be stored (DQ-1).
var ErrNotRegistered = errors.New("source is not in the register")

// Record is one word-meaning pair read from a source file.
type Record struct {
	Text      string // the word or phrase in the source language
	English   string
	Reference string // where in the file
}

// Parser reads a source file into records.
type Parser interface {
	Parse(r io.Reader) iter.Seq2[Record, error]
}

// Feed is one importable file of a registered source.
type Feed struct {
	Name     string // command-line name, e.g. "gatitos-bem"
	Source   string // register name, e.g. "GATITOS"
	Language string
	URL      string // pinned to a revision, so the bytes never change
	SHA256   string
	Core     bool // its English glosses seed the core concept list
	Parser   Parser
}

// Stats counts what an ingest created.
type Stats struct {
	Records, Entries, Senses int
}

// Beginner starts a transaction. A pool or an outer transaction both work.
type Beginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// Ingest stores every record of a feed as `sourced` lexicon data, in one
// transaction, with the audit actor set to the feed.
func Ingest(ctx context.Context, db Beginner, feed Feed, r io.Reader) (Stats, error) {
	var stats Stats
	tx, err := db.Begin(ctx)
	if err != nil {
		return stats, fmt.Errorf("begin: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	src, err := dbgen.New(tx).GetSourceByName(ctx, feed.Source)
	if errors.Is(err, pgx.ErrNoRows) {
		return stats, fmt.Errorf("%s: %w", feed.Source, ErrNotRegistered)
	}
	if err != nil {
		return stats, fmt.Errorf("get source: %w", err)
	}
	if _, err := tx.Exec(ctx, "SELECT set_config('app.actor_id', $1, true)", "ingest:"+feed.Name); err != nil {
		return stats, fmt.Errorf("set actor: %w", err)
	}

	lex := lexicon.NewStore(tx)
	for rec, err := range feed.Parser.Parse(r) {
		if err != nil {
			return stats, err
		}
		added, err := lex.AddSourced(ctx, lexicon.SourcedSense{
			Language:  feed.Language,
			Headword:  rec.Text,
			GlossEn:   rec.English,
			Core:      feed.Core,
			SourceID:  src.ID,
			Reference: rec.Reference,
		})
		if err != nil {
			return stats, err
		}
		stats.Records++
		if added.Entry {
			stats.Entries++
		}
		if added.Sense {
			stats.Senses++
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return stats, fmt.Errorf("commit: %w", err)
	}
	return stats, nil
}
