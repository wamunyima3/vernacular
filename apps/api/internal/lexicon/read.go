package lexicon

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/textnorm"
)

// Sense is a servable meaning of an entry.
type Sense struct {
	ID         uuid.UUID
	Gloss      string
	Definition string
	ConceptID  uuid.NullUUID
	Examples   []Example // filled by Get, not by Search
	Trust      items.Trust
}

// Example is a servable example sentence of a sense (LEX-5).
type Example struct {
	ID                  uuid.UUID
	Text, TranslationEn string
	Trust               items.Trust
}

// Variant is another accepted or seen spelling of an entry (LEX-3).
type Variant struct{ Spelling, Note string }

// Entry is a dictionary entry as the read endpoints return it (LEX-1, LEX-3,
// LEX-6). Optional fields are empty when unknown.
type Entry struct {
	ID                                                  uuid.UUID
	Language, Headword, PartOfSpeech, NounClass, Plural string
	Root, Variety, Register                             string
	Variants                                            []Variant // filled by Get, not by Search
	Senses                                              []Sense
	Trust                                               items.Trust
}

// Hit is an entry found by search and what the query matched (US-05).
type Hit struct {
	Match string // headword, variant, root or similar
	Entry Entry
}

// Search finds a language's entries for a query, best matches first (LEX-4).
// The query is folded (NFR-7), spelled as the style guide does, and stemmed,
// so variants and prefixed or inflected forms match.
func (s *Store) Search(ctx context.Context, language, query string, limit int) ([]Hit, error) {
	fold := textnorm.Fold(query)
	standard, _ := Standardise(language, fold)
	rows, err := s.q.SearchEntries(ctx, dbgen.SearchEntriesParams{
		Language: language, Fold: fold, Standard: standard, Stem: Stem(language, standard), Lim: int32(limit), //nolint:gosec // limit is 1..50, checked by the contract
	})
	if err != nil {
		return nil, fmt.Errorf("search entries: %w", err)
	}
	ids := make([]uuid.UUID, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
	}
	entries, err := s.entries(ctx, ids)
	if err != nil {
		return nil, err
	}
	hits := make([]Hit, 0, len(rows))
	for _, r := range rows {
		if e, ok := entries[r.ID]; ok {
			hits = append(hits, Hit{Match: r.Matched, Entry: e})
		}
	}
	return hits, nil
}

// Get reads a full servable entry: variants, senses and their examples
// (US-16). False when there is no such entry or it cannot be served.
func (s *Store) Get(ctx context.Context, id uuid.UUID) (Entry, bool, error) {
	entries, err := s.entries(ctx, []uuid.UUID{id})
	e, ok := entries[id]
	if err != nil || !ok {
		return Entry{}, false, err
	}
	variants, err := s.q.ListVariants(ctx, id)
	if err != nil {
		return Entry{}, false, fmt.Errorf("list variants: %w", err)
	}
	for _, v := range variants {
		e.Variants = append(e.Variants, Variant{Spelling: v.Spelling, Note: v.Note})
	}
	if err := s.withExamples(ctx, e.Senses); err != nil {
		return Entry{}, false, err
	}
	return e, true, nil
}

// entries reads servable entries with their servable senses and everyone's
// trust fields, in four queries however many there are.
func (s *Store) entries(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]Entry, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := s.q.ListEntries(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("list entries: %w", err)
	}
	senses, err := s.q.ListServableSenses(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("list senses: %w", err)
	}
	itemIDs := slices.Clone(ids)
	for _, sn := range senses {
		itemIDs = append(itemIDs, sn.ID)
	}
	trusts, err := s.items.Trusts(ctx, itemIDs)
	if err != nil {
		return nil, err
	}
	byEntry := make(map[uuid.UUID][]Sense, len(rows))
	for _, sn := range senses {
		byEntry[sn.EntryID] = append(byEntry[sn.EntryID], Sense{
			ID: sn.ID, Gloss: sn.Gloss, Definition: sn.Definition, ConceptID: sn.ConceptID, Trust: trusts[sn.ID],
		})
	}
	out := make(map[uuid.UUID]Entry, len(rows))
	for _, r := range rows {
		if len(byEntry[r.ID]) == 0 {
			continue // nothing about it can be served
		}
		out[r.ID] = Entry{
			ID: r.ID, Language: r.LanguageCode, Headword: r.Headword, PartOfSpeech: string(r.PartOfSpeech.PartOfSpeech),
			NounClass: r.NounClass.String, Plural: r.Plural.String, Root: r.Root.String, Variety: r.Variety.String,
			Register: string(r.Register), Senses: byEntry[r.ID], Trust: trusts[r.ID],
		}
	}
	return out, nil
}

// withExamples fills in the senses' servable examples with their trust.
func (s *Store) withExamples(ctx context.Context, senses []Sense) error {
	ids := make([]uuid.UUID, len(senses))
	for i, sn := range senses {
		ids[i] = sn.ID
	}
	rows, err := s.q.ListServableExamples(ctx, ids)
	if err != nil || len(rows) == 0 {
		return err
	}
	exIDs := make([]uuid.UUID, len(rows))
	for i, r := range rows {
		exIDs[i] = r.ID
	}
	trusts, err := s.items.Trusts(ctx, exIDs)
	if err != nil {
		return err
	}
	bySense := make(map[uuid.UUID][]Example, len(senses))
	for _, r := range rows {
		bySense[r.SenseID] = append(bySense[r.SenseID], Example{ID: r.ID, Text: r.Text, TranslationEn: r.TranslationEn, Trust: trusts[r.ID]})
	}
	for i := range senses {
		senses[i].Examples = bySense[senses[i].ID]
	}
	return nil
}

// ConceptEntries is a concept with the servable entries that express it in
// any language (LEX-7). Each entry lists only the senses linked to the concept.
type ConceptEntries struct {
	Concept
	Entries []Entry
}

// ConceptEntries reads a concept across all languages. False when there is none.
func (s *Store) ConceptEntries(ctx context.Context, id uuid.UUID) (ConceptEntries, bool, error) {
	c, err := s.Concept(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return ConceptEntries{}, false, nil
	}
	if err != nil {
		return ConceptEntries{}, false, err
	}
	rows, err := s.q.ListConceptEntries(ctx, uuid.NullUUID{UUID: id, Valid: true})
	if err != nil {
		return ConceptEntries{}, false, fmt.Errorf("list concept entries: %w", err)
	}
	ids := make([]uuid.UUID, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
	}
	entries, err := s.entries(ctx, ids)
	if err != nil {
		return ConceptEntries{}, false, err
	}
	out := ConceptEntries{Concept: c}
	for _, eid := range ids {
		e := entries[eid]
		e.Senses = slices.DeleteFunc(e.Senses, func(sn Sense) bool { return sn.ConceptID.UUID != id })
		if len(e.Senses) > 0 {
			out.Entries = append(out.Entries, e)
		}
	}
	return out, true, nil
}
