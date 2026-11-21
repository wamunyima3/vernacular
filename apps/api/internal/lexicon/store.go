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
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/textnorm"
)

// ErrInvalid wraps errors in what a caller proposed, as opposed to failures.
var ErrInvalid = errors.New("invalid")

// Proposal is one word-meaning pair to add, from a source or a model.
type Proposal struct {
	Language     string
	Headword     string
	PartOfSpeech string // empty when unknown
	NounClass    string // "1/2" style, checked by the table; set only on new entries
	Plural       string // set only on new entries
	GlossEn      string
	Core         bool          // the concept belongs to the core list (LEX-7)
	ConceptID    uuid.NullUUID // when known; otherwise found or created by gloss
	Origin       items.Origin
}

// Added reports which rows an Add call created. When Sense is false, the
// sense already existed and the origin was added to its sources.
type Added struct {
	Entry, Sense, Variant bool
	EntryID, SenseID      uuid.UUID
}

// Store reads and writes lexicon tables.
type Store struct {
	q      *dbgen.Queries
	items  *items.Store
	locked map[string]bool // languages this transaction holds the write lock for
}

// NewStore wraps a pool or transaction. Writes that span several rows should
// be given a transaction.
func NewStore(db dbgen.DBTX) *Store {
	return &Store{q: dbgen.New(db), items: items.NewStore(db), locked: map[string]bool{}}
}

// Add normalises the headword and gloss (textnorm), then records the sense
// with the proposal's origin, creating its concept and entry when they do not
// exist yet. When the entry or sense already exists, its status is kept and
// the origin is added to its sources, so agreement between sources and models
// is visible (AI-3). Repeating a call changes nothing, so imports and
// candidate runs can be re-run.
func (s *Store) Add(ctx context.Context, p Proposal) (Added, error) {
	var added Added
	p.Headword, p.GlossEn = textnorm.Headword(p.Headword), textnorm.Gloss(p.GlossEn)
	if p.Headword == "" || p.GlossEn == "" {
		return added, fmt.Errorf("%w: headword and English gloss are required", ErrInvalid)
	}
	conceptID := p.ConceptID
	if !conceptID.Valid {
		concept, err := s.q.UpsertConcept(ctx, dbgen.UpsertConceptParams{GlossEn: p.GlossEn, Core: p.Core, Domain: DomainGeneral})
		if err != nil {
			return added, fmt.Errorf("upsert concept %q: %w", p.GlossEn, err)
		}
		conceptID = uuid.NullUUID{UUID: concept.ID, Valid: true}
	}

	pos := dbgen.NullPartOfSpeech{PartOfSpeech: dbgen.PartOfSpeech(p.PartOfSpeech), Valid: p.PartOfSpeech != ""}
	// Spellings that differ only by orthography merge into one entry under the
	// style guide's standard form; the proposed spelling becomes a variant.
	headword, variant := Standardise(p.Language, p.Headword)
	if err := s.lock(ctx, p.Language); err != nil {
		return added, err
	}
	entryID, err := s.q.FindEntry(ctx, dbgen.FindEntryParams{LanguageCode: p.Language, Headword: headword, PartOfSpeech: pos})
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		added.Entry = true
		entryID, err = s.items.Create(ctx, items.KindEntry, p.Language, p.Origin, func(id uuid.UUID) error {
			return s.q.CreateEntry(ctx, dbgen.CreateEntryParams{
				ID: id, LanguageCode: p.Language, Headword: headword, PartOfSpeech: pos,
				Stem:      pgtype.Text{String: Stem(p.Language, textnorm.Fold(headword)), Valid: true},
				NounClass: database.Text(strings.TrimSpace(p.NounClass)), Plural: database.Text(textnorm.Headword(p.Plural)),
			})
		})
	case err == nil:
		err = s.items.Link(ctx, entryID, p.Origin)
	}
	if err != nil {
		return added, fmt.Errorf("entry %q: %w", headword, err)
	}

	if variant != "" {
		n, err := s.q.AddVariant(ctx, dbgen.AddVariantParams{EntryID: entryID, Spelling: variant, Note: "as written in " + p.Origin.Reference})
		if err != nil {
			return added, fmt.Errorf("variant %q: %w", variant, err)
		}
		added.Variant = n > 0
	}

	added.EntryID = entryID
	senseID, err := s.q.FindSense(ctx, dbgen.FindSenseParams{EntryID: entryID, ConceptID: conceptID})
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		added.Sense = true
		senseID, err = s.items.Create(ctx, items.KindSense, p.Language, p.Origin, func(id uuid.UUID) error {
			return s.q.CreateSense(ctx, dbgen.CreateSenseParams{
				ID: id, EntryID: entryID, ConceptID: conceptID, Gloss: p.GlossEn,
			})
		})
	case err == nil:
		err = s.items.Link(ctx, senseID, p.Origin)
	}
	if err != nil {
		return added, fmt.Errorf("sense %q = %q: %w", headword, p.GlossEn, err)
	}
	added.SenseID = senseID
	return added, nil
}

// AddExample records a usage example of a sense with its English
// translation, as a reviewable item of its own.
func (s *Store) AddExample(ctx context.Context, language string, senseID uuid.UUID, text, translationEn string, origin items.Origin) (uuid.UUID, error) {
	text, translationEn = textnorm.Squash(textnorm.NFC(text)), textnorm.Squash(textnorm.NFC(translationEn))
	if text == "" || translationEn == "" {
		return uuid.Nil, fmt.Errorf("%w: example text and translation are required", ErrInvalid)
	}
	id, err := s.items.Create(ctx, items.KindExample, language, origin, func(id uuid.UUID) error {
		return s.q.CreateExample(ctx, dbgen.CreateExampleParams{ID: id, SenseID: senseID, Text: text, TranslationEn: translationEn})
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("example: %w", err)
	}
	return id, nil
}

// lock takes the language's write lock once per transaction, so concurrent
// writers (parallel AI jobs, an import) cannot both create the same entry.
// One key per language: no lock ordering to get wrong, and a long import
// holds one lock rather than thousands. Released at commit.
func (s *Store) lock(ctx context.Context, language string) error {
	if s.locked[language] {
		return nil
	}
	if err := s.q.LockLanguage(ctx, language); err != nil {
		return fmt.Errorf("lock %s: %w", language, err)
	}
	s.locked[language] = true
	return nil
}

// Restem recomputes the stem of every entry in a language, after the
// stemmer's rules change. It returns how many stems changed.
func (s *Store) Restem(ctx context.Context, language string) (int, error) {
	rows, err := s.q.ListHeadwords(ctx, language)
	if err != nil {
		return 0, fmt.Errorf("list headwords: %w", err)
	}
	changed := 0
	for _, r := range rows {
		stem := Stem(language, textnorm.Fold(r.Headword))
		if r.Stem.Valid && r.Stem.String == stem {
			continue
		}
		if err := s.q.SetStem(ctx, dbgen.SetStemParams{ID: r.ID, Stem: pgtype.Text{String: stem, Valid: true}}); err != nil {
			return changed, fmt.Errorf("set stem of %q: %w", r.Headword, err)
		}
		changed++
	}
	return changed, nil
}

// DomainGeneral is the domain of concepts no list has classified yet.
const DomainGeneral = "general"

// AddCoreConcept adds an English concept to the core list (LEX-7) under a
// domain such as "health". It reports whether the concept is new.
func (s *Store) AddCoreConcept(ctx context.Context, glossEn, domain string) (bool, error) {
	domain = strings.ToLower(textnorm.Squash(domain))
	if domain == "" {
		return false, fmt.Errorf("%w: domain is required", ErrInvalid)
	}
	row, err := s.upsertConcept(ctx, glossEn, true, domain)
	return row.Inserted, err
}

// ConceptFor finds or creates the general concept for an English gloss,
// without adding it to the core list, and returns its ID.
func (s *Store) ConceptFor(ctx context.Context, glossEn string) (uuid.UUID, error) {
	row, err := s.upsertConcept(ctx, glossEn, false, DomainGeneral)
	return row.ID, err
}

func (s *Store) upsertConcept(ctx context.Context, glossEn string, core bool, domain string) (dbgen.UpsertConceptRow, error) {
	glossEn = textnorm.Gloss(glossEn)
	if glossEn == "" {
		return dbgen.UpsertConceptRow{}, fmt.Errorf("%w: gloss is required", ErrInvalid)
	}
	row, err := s.q.UpsertConcept(ctx, dbgen.UpsertConceptParams{GlossEn: glossEn, Core: core, Domain: domain})
	if err != nil {
		return dbgen.UpsertConceptRow{}, fmt.Errorf("upsert concept %q: %w", glossEn, err)
	}
	return row, nil
}

// Concept is an English meaning that entries in every language point to.
type Concept struct {
	ID         uuid.UUID
	GlossEn    string
	Definition string
	Domain     string
	Core       bool // part of the core list (LEX-7)
}

// Concept returns one concept.
func (s *Store) Concept(ctx context.Context, id uuid.UUID) (Concept, error) {
	c, err := s.q.GetConcept(ctx, id)
	if err != nil {
		return Concept{}, fmt.Errorf("get concept %s: %w", id, err)
	}
	return Concept{ID: c.ID, GlossEn: c.GlossEn, Definition: c.Definition, Domain: c.Domain, Core: c.Core}, nil
}

// CoreConcepts lists the IDs of the core list (LEX-7).
func (s *Store) CoreConcepts(ctx context.Context) ([]uuid.UUID, error) {
	ids, err := s.q.ListCoreConcepts(ctx)
	if err != nil {
		return nil, fmt.Errorf("list core concepts: %w", err)
	}
	return ids, nil
}

// PartsOfSpeech lists the values the part_of_speech column accepts.
func PartsOfSpeech() []string {
	all := dbgen.AllPartOfSpeechValues()
	out := make([]string, len(all))
	for i, p := range all {
		out[i] = string(p)
	}
	return out
}

// SenseCard is a sense as a reviewer sees it: the word, without its sources.
type SenseCard struct {
	ID                                        uuid.UUID
	ConceptID                                 uuid.NullUUID
	GlossEn                                   string // the concept's, if any
	Headword, PartOfSpeech, NounClass, Plural string
}

// SenseCards reads the given senses, in the order of ids.
func (s *Store) SenseCards(ctx context.Context, ids []uuid.UUID) ([]SenseCard, error) {
	rows, err := s.q.ListSenseCards(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("list sense cards: %w", err)
	}
	out := make([]SenseCard, len(rows))
	for i, r := range rows {
		out[i] = SenseCard{
			ID: r.ID, ConceptID: r.ConceptID, GlossEn: r.GlossEn.String, Headword: r.Headword, PartOfSpeech: string(r.PartOfSpeech.PartOfSpeech),
			NounClass: r.NounClass.String, Plural: r.Plural.String,
		}
	}
	return out, nil
}

// EntryCount is how many entries a language has, verified and in all.
type EntryCount struct{ Verified, Total int }

// CountEntries counts each language's entries; languages without any are left out.
func (s *Store) CountEntries(ctx context.Context) (map[string]EntryCount, error) {
	rows, err := s.q.CountEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("count entries: %w", err)
	}
	out := make(map[string]EntryCount, len(rows))
	for _, r := range rows {
		out[r.LanguageCode] = EntryCount{Verified: int(r.Verified), Total: int(r.Total)}
	}
	return out, nil
}
