package lexicon

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/textnorm"
)

// Sense is a servable meaning of an entry.
type Sense struct {
	ID    uuid.UUID
	Gloss string
	Trust items.Trust
}

// Entry is a dictionary entry as the read endpoints return it (LEX-1, LEX-3,
// LEX-6). Optional fields are empty when unknown.
type Entry struct {
	ID                                                  uuid.UUID
	Language, Headword, PartOfSpeech, NounClass, Plural string
	Root, Variety, Register                             string
	Trust                                               items.Trust
	Senses                                              []Sense
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
	entries := make([]Entry, len(rows))
	for i, r := range rows {
		entries[i] = Entry{
			ID: r.ID, Language: r.LanguageCode, Headword: r.Headword, PartOfSpeech: string(r.PartOfSpeech.PartOfSpeech),
			NounClass: r.NounClass.String, Plural: r.Plural.String, Root: r.Root.String, Variety: r.Variety.String, Register: string(r.Register),
		}
	}
	if err := s.withSenses(ctx, entries); err != nil {
		return nil, err
	}
	hits := make([]Hit, len(rows))
	for i, r := range rows {
		hits[i] = Hit{Match: r.Matched, Entry: entries[i]}
	}
	return hits, nil
}

// withSenses fills in the entries' servable senses and everyone's trust
// fields, in three queries however many entries there are.
func (s *Store) withSenses(ctx context.Context, entries []Entry) error {
	if len(entries) == 0 {
		return nil
	}
	ids := make([]uuid.UUID, len(entries))
	for i, e := range entries {
		ids[i] = e.ID
	}
	senses, err := s.q.ListServableSenses(ctx, ids)
	if err != nil {
		return fmt.Errorf("list senses: %w", err)
	}
	for _, sn := range senses {
		ids = append(ids, sn.ID)
	}
	trusts, err := s.items.Trusts(ctx, ids)
	if err != nil {
		return err
	}
	byEntry := make(map[uuid.UUID][]Sense, len(entries))
	for _, sn := range senses {
		byEntry[sn.EntryID] = append(byEntry[sn.EntryID], Sense{ID: sn.ID, Gloss: sn.Gloss, Trust: trusts[sn.ID]})
	}
	for i := range entries {
		entries[i].Trust, entries[i].Senses = trusts[entries[i].ID], byEntry[entries[i].ID]
	}
	return nil
}
