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
	AddSourced(ctx context.Context, in lexicon.SourcedSense) (lexicon.Added, error)
	AddCoreConcept(ctx context.Context, glossEn, domain string) (bool, error)
}

// newLexicon gives a writer bound to the import transaction.
var newLexicon = func(tx pgx.Tx) LexiconWriter { return lexicon.NewStore(tx) }

// Stats counts what an ingest created.
type Stats struct {
	Records, Entries, Senses, Variants int
}

// Ingest stores every record of a feed as `sourced` lexicon data, in one
// transaction, with the audit actor set to the feed.
func Ingest(ctx context.Context, db database.Beginner, feed Feed, r io.Reader) (Stats, error) {
	var stats Stats
	err := database.InTx(ctx, db, "ingest:"+feed.Name, func(tx pgx.Tx) error {
		src, err := dbgen.New(tx).GetSourceByName(ctx, feed.Source)
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", feed.Source, ErrNotRegistered)
		}
		if err != nil {
			return fmt.Errorf("get source: %w", err)
		}
		lex := newLexicon(tx)
		for rec, err := range feed.Parser.Parse(r) {
			if err != nil {
				return err
			}
			added, err := lex.AddSourced(ctx, lexicon.SourcedSense{
				Language: feed.Language, Headword: rec.Text, GlossEn: rec.English, Core: feed.Core,
				SourceID: src.ID, Reference: rec.Reference,
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
