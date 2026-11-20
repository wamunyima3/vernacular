// Package lexicon owns concepts, entries, variants, senses and examples
// (architecture §5). Other modules add or read lexicon data only through it.
package lexicon

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/textnorm"
)

// SourcedSense is one word-meaning pair taken from a licensed source.
type SourcedSense struct {
	Language     string
	Headword     string
	PartOfSpeech string // empty when the source does not say
	GlossEn      string
	Core         bool // the concept belongs to the core list (LEX-7)
	SourceID     uuid.UUID
	Reference    string // where in the source, e.g. "en_bem.jsonl:12"
}

// Added reports which rows an AddSourced call created.
type Added struct {
	Entry, Sense bool
}

// Store reads and writes lexicon tables.
type Store struct {
	q *dbgen.Queries
}

// NewStore wraps a pool or transaction. Writes that span several rows should
// be given a transaction.
func NewStore(db dbgen.DBTX) *Store {
	return &Store{q: dbgen.New(db)}
}

// AddSourced normalises the headword and gloss (textnorm), then records a sense with status `sourced`, creating its concept and
// entry when they do not exist yet. Repeating a call changes nothing, so
// imports can be re-run.
func (s *Store) AddSourced(ctx context.Context, in SourcedSense) (Added, error) {
	var added Added
	in.Headword, in.GlossEn = textnorm.Headword(in.Headword), textnorm.Gloss(in.GlossEn)
	if in.Headword == "" || in.GlossEn == "" {
		return added, errors.New("headword and English gloss are required")
	}
	concept, err := s.q.UpsertConcept(ctx, dbgen.UpsertConceptParams{GlossEn: in.GlossEn, Core: in.Core, Domain: DomainGeneral})
	if err != nil {
		return added, fmt.Errorf("upsert concept %q: %w", in.GlossEn, err)
	}
	conceptID := uuid.NullUUID{UUID: concept.ID, Valid: true}

	pos := dbgen.NullPartOfSpeech{PartOfSpeech: dbgen.PartOfSpeech(in.PartOfSpeech), Valid: in.PartOfSpeech != ""}
	headword := in.Headword
	entryID, err := s.q.FindEntry(ctx, dbgen.FindEntryParams{LanguageCode: in.Language, Headword: headword, PartOfSpeech: pos})
	if errors.Is(err, pgx.ErrNoRows) {
		added.Entry = true
		entryID, err = s.newItem(ctx, dbgen.ItemKindEntry, in, func(id uuid.UUID) error {
			return s.q.CreateEntry(ctx, dbgen.CreateEntryParams{ID: id, LanguageCode: in.Language, Headword: headword, PartOfSpeech: pos})
		})
	}
	if err != nil {
		return added, fmt.Errorf("entry %q: %w", headword, err)
	}

	_, err = s.q.FindSense(ctx, dbgen.FindSenseParams{EntryID: entryID, ConceptID: conceptID})
	if errors.Is(err, pgx.ErrNoRows) {
		added.Sense = true
		_, err = s.newItem(ctx, dbgen.ItemKindSense, in, func(id uuid.UUID) error {
			return s.q.CreateSense(ctx, dbgen.CreateSenseParams{
				ID: id, EntryID: entryID, ConceptID: conceptID, Gloss: in.GlossEn,
			})
		})
	}
	if err != nil {
		return added, fmt.Errorf("sense %q = %q: %w", headword, in.GlossEn, err)
	}
	return added, nil
}

// DomainGeneral is the domain of concepts no list has classified yet.
const DomainGeneral = "general"

// AddCoreConcept adds an English concept to the core list (LEX-7) under a
// domain such as "health". It reports whether the concept is new.
func (s *Store) AddCoreConcept(ctx context.Context, glossEn, domain string) (bool, error) {
	glossEn, domain = textnorm.Gloss(glossEn), strings.ToLower(textnorm.Squash(domain))
	if glossEn == "" || domain == "" {
		return false, errors.New("gloss and domain are required")
	}
	row, err := s.q.UpsertConcept(ctx, dbgen.UpsertConceptParams{GlossEn: glossEn, Core: true, Domain: domain})
	if err != nil {
		return false, fmt.Errorf("upsert concept %q: %w", glossEn, err)
	}
	return row.Inserted, nil
}

// newItem creates the shared items row (ADR-0006), then the typed row through
// create, then links the source. It is the only way lexicon rows are created.
func (s *Store) newItem(ctx context.Context, kind dbgen.ItemKind, in SourcedSense, create func(uuid.UUID) error) (uuid.UUID, error) {
	id, err := s.q.CreateItem(ctx, dbgen.CreateItemParams{Kind: kind, LanguageCode: in.Language, Status: dbgen.TrustStatusSourced})
	if err != nil {
		return id, fmt.Errorf("create item: %w", err)
	}
	if err := create(id); err != nil {
		return id, err
	}
	err = s.q.LinkSource(ctx, dbgen.LinkSourceParams{ItemID: id, SourceID: in.SourceID, Reference: in.Reference})
	if err != nil {
		return id, fmt.Errorf("link source: %w", err)
	}
	return id, nil
}
