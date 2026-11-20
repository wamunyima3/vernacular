// Package languages owns the list of supported and planned languages.
package languages

import "context"

// Language is the domain view of a language.
type Language struct {
	Code    string
	Name    string
	OwnName string
	Wave    int
	Status  string
}

type store interface {
	List(ctx context.Context) ([]Language, error)
}

// Service holds the language rules. Entry counts (P1-API-01) will be added here.
type Service struct {
	store store
}

// NewService returns a Service backed by the given store.
func NewService(s store) *Service {
	return &Service{store: s}
}

// List returns all languages in rollout order.
func (s *Service) List(ctx context.Context) ([]Language, error) {
	return s.store.List(ctx)
}
