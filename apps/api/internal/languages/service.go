// Package languages owns the list of supported and planned languages.
package languages

import (
	"context"

	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
)

// Language is the domain view of a language.
type Language struct {
	Code    string
	Name    string
	OwnName string
	Wave    int
	Status  string
	Entries lexicon.EntryCount
}

type store interface {
	List(ctx context.Context) ([]Language, error)
}

// Counter counts entries per language; the lexicon module owns entries.
type Counter interface {
	CountEntries(ctx context.Context) (map[string]lexicon.EntryCount, error)
}

// Service holds the language rules.
type Service struct {
	store   store
	counter Counter
}

// NewService returns a Service backed by the given store and entry counter.
func NewService(s store, c Counter) *Service {
	return &Service{store: s, counter: c}
}

// List returns all languages in rollout order, with live entry counts.
func (s *Service) List(ctx context.Context) ([]Language, error) {
	langs, err := s.store.List(ctx)
	if err != nil {
		return nil, err
	}
	counts, err := s.counter.CountEntries(ctx)
	if err != nil {
		return nil, err
	}
	for i := range langs {
		langs[i].Entries = counts[langs[i].Code]
	}
	return langs, nil
}
