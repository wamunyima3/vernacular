package lexicon

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/textnorm"
)

// toEnglishEntries is how many entries ToEnglish searches to fill its
// concepts; several entries often share one.
const toEnglishEntries = 50

// LookupHit is a concept a lookup found, with the entries expressing it and
// what the query matched.
type LookupHit struct {
	Match     string
	ConceptID uuid.UUID
	GlossEn   string
	Entries   []Entry
}

// FromEnglish finds concepts by English gloss and their entries in the
// languages, all when none are given (LEX-2, US-04). Concepts with no
// servable entry there are left out.
func (s *Store) FromEnglish(ctx context.Context, query string, languages []string, limit int, f Filter) ([]LookupHit, error) {
	concepts, err := s.q.FindConcepts(ctx, dbgen.FindConceptsParams{Fold: textnorm.Fold(query), Lim: int32(limit)}) //nolint:gosec // limit is 1..20, checked by the contract
	if err != nil {
		return nil, fmt.Errorf("find concepts: %w", err)
	}
	ids := database.IDs(concepts, func(c dbgen.FindConceptsRow) uuid.UUID { return c.ID })
	byConcept, err := s.conceptEntries(ctx, ids, languages, f)
	if err != nil {
		return nil, err
	}
	var hits []LookupHit
	for _, c := range concepts {
		if entries := byConcept[c.ID]; len(entries) > 0 {
			match := MatchSimilar
			if c.Exact {
				match = MatchHeadword
			}
			hits = append(hits, LookupHit{Match: match, ConceptID: c.ID, GlossEn: c.GlossEn, Entries: entries})
		}
	}
	return hits, nil
}

// ToEnglish searches a language's entries and gives each meaning found in
// English: one hit per concept, in the order search ranked them.
func (s *Store) ToEnglish(ctx context.Context, language, query string, limit int, f Filter) ([]LookupHit, error) {
	found, err := s.Search(ctx, language, query, toEnglishEntries, f)
	if err != nil {
		return nil, err
	}
	var hits []LookupHit
	index := map[uuid.UUID]int{}
	for _, h := range found {
		for _, c := range conceptsOf(h.Entry) {
			if i, ok := index[c]; ok {
				hits[i].Entries = append(hits[i].Entries, mustOnly(h.Entry, c))
			} else if len(hits) < limit {
				index[c] = len(hits)
				hits = append(hits, LookupHit{Match: h.Match, ConceptID: c, Entries: []Entry{mustOnly(h.Entry, c)}})
			}
		}
	}
	if len(hits) == 0 {
		return nil, nil
	}
	ids := database.IDs(hits, func(h LookupHit) uuid.UUID { return h.ConceptID })
	concepts, err := s.q.ListConcepts(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("list concepts: %w", err)
	}
	for _, c := range concepts {
		hits[index[c.ID]].GlossEn = c.GlossEn
	}
	return hits, nil
}

// conceptsOf lists the distinct concepts of an entry's senses, in order.
func conceptsOf(e Entry) []uuid.UUID {
	var out []uuid.UUID
	for _, sn := range e.Senses {
		if sn.ConceptID.Valid && !slices.Contains(out, sn.ConceptID.UUID) {
			out = append(out, sn.ConceptID.UUID)
		}
	}
	return out
}

// mustOnly is onlyConcept for a concept known to be among the entry's.
func mustOnly(e Entry, concept uuid.UUID) Entry {
	e, _ = onlyConcept(e, concept)
	return e
}
