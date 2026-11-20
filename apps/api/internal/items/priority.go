package items

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/wamunyima3/vernacular/apps/api/internal/items/dbgen"
)

// Review priorities (AI-3, US-21).
const (
	PriorityQuick    = dbgen.ReviewPriorityQuick    // every source that answered agrees
	PriorityStandard = dbgen.ReviewPriorityStandard // at least two sources agree
	PriorityFull     = dbgen.ReviewPriorityFull     // a single source, or disagreement
)

// Score is an item's agreement and the priority it earns.
type Score struct {
	Agreement int
	Priority  dbgen.ReviewPriority
}

// Scores rates rival items: different answers to one question, such as the
// senses of one concept in one language. Each item maps to the distinct
// sources that attest it. An answer every responding source gives is a
// quick confirm; one that two or more give is a standard review; anything
// else needs a full review.
func Scores(rivals map[uuid.UUID][]uuid.UUID) map[uuid.UUID]Score {
	answered := map[uuid.UUID]bool{}
	for _, sources := range rivals {
		for _, s := range sources {
			answered[s] = true
		}
	}
	scores := make(map[uuid.UUID]Score, len(rivals))
	for id, sources := range rivals {
		n := len(sources)
		priority := PriorityFull
		switch {
		case n >= 2 && n == len(answered):
			priority = PriorityQuick
		case n >= 2:
			priority = PriorityStandard
		}
		scores[id] = Score{Agreement: n, Priority: priority}
	}
	return scores
}

// Prioritise scores rival items from their sources and stores the result.
func (s *Store) Prioritise(ctx context.Context, rivals []uuid.UUID) error {
	if len(rivals) == 0 {
		return nil
	}
	rows, err := s.q.ListItemSources(ctx, rivals)
	if err != nil {
		return fmt.Errorf("list item sources: %w", err)
	}
	bySource := make(map[uuid.UUID][]uuid.UUID, len(rivals))
	for _, id := range rivals {
		bySource[id] = nil
	}
	for _, r := range rows {
		bySource[r.ItemID] = append(bySource[r.ItemID], r.SourceID)
	}
	for id, sc := range Scores(bySource) {
		err := s.q.SetPriority(ctx, dbgen.SetPriorityParams{ID: id, ReviewPriority: sc.Priority, Agreement: int16(sc.Agreement)})
		if err != nil {
			return fmt.Errorf("set priority of %s: %w", id, err)
		}
	}
	return nil
}
