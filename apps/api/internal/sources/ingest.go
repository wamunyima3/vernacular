// Package sources owns the source register (DQ-1) and ingestion: one parser
// per licensed source, all feeding the lexicon through the same path.
package sources

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
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
	Name     string // command-line name, e.g. "gatitos-en-bem"
	Source   string // register name, e.g. "GATITOS"
	Language string
	URL      string // pinned to a revision, so the bytes never change
	SHA256   string
	Core     bool // its English glosses seed the core concept list
	Parser   Parser
}

// LexiconWriter is what ingestion needs from the lexicon module.
type LexiconWriter interface {
	Add(ctx context.Context, p lexicon.Proposal) (lexicon.Added, error)
	AddCoreConcept(ctx context.Context, glossEn, domain string) (bool, error)
}

// newLexicon gives a writer bound to the import transaction.
var newLexicon = func(tx pgx.Tx) LexiconWriter { return lexicon.NewStore(tx) }

// Lookup returns the ID of a registered source. Anything stored from a source
// must cite it, so callers look it up first (DQ-1).
func Lookup(ctx context.Context, db dbgen.DBTX, name string) (uuid.UUID, error) {
	src, err := dbgen.New(db).GetSourceByName(ctx, name)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("%s: %w", name, ErrNotRegistered)
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("get source %s: %w", name, err)
	}
	return src.ID, nil
}

// Stats counts what an ingest created.
type Stats struct {
	Records, Entries, Senses, Variants int
}

// Ingest stores every record of a feed as `sourced` lexicon data, in one
// transaction, with the audit actor set to the feed.
func Ingest(ctx context.Context, db database.Beginner, feed Feed, r io.Reader) (Stats, error) {
	var stats Stats
	err := database.InTx(ctx, db, "ingest:"+feed.Name, func(tx pgx.Tx) error {
		sourceID, err := Lookup(ctx, tx, feed.Source)
		if err != nil {
			return err
		}
		lex := newLexicon(tx)
		for rec, err := range feed.Parser.Parse(r) {
			if err != nil {
				return err
			}
			added, err := lex.Add(ctx, lexicon.Proposal{
				Language: feed.Language, Headword: rec.Text, GlossEn: rec.English, Core: feed.Core,
				Origin: items.Sourced(sourceID, rec.Reference),
			})
			if err != nil {
				return err
			}
			stats.Records++
			if added.Entry {
				stats.Entries++
			}
			if added.Sense {
				stats.Senses++
			}
			if added.Variant {
				stats.Variants++
			}
		}
		return nil
	})
	if err != nil {
		return Stats{}, err
	}
	return stats, nil
}
