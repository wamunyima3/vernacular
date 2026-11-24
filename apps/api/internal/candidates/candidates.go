// Package candidates asks language models for words that express a concept
// (AI-1). Everything they propose is stored as `machine`, citing the
// provider's source row and the exact model version (AI-2), and waits for
// human review. Candidates are never used as training data (DQ-10).
package candidates

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/phrases"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
)

// Request is one concept to express in one language.
type Request struct {
	Language string
	Concept  lexicon.Concept
}

// Candidate is one word a model proposed.
type Candidate struct {
	Headword     string `json:"headword"`
	PartOfSpeech string `json:"part_of_speech"`
	NounClass    string `json:"noun_class"`
	Plural       string `json:"plural"`
}

// Result is what a model returned, and which model version answered.
type Result struct {
	Model      string
	Candidates []Candidate
}

// PhraseRequest is one line of the phrase source list to say in one language.
type PhraseRequest struct {
	Language, Pack, SourceEn, Notes string
}

// PhraseCandidate is one way a model proposed to say a phrase.
type PhraseCandidate struct {
	Text  string `json:"text"`
	Notes string `json:"notes"`
}

// PhraseResult is what a model returned for a phrase.
type PhraseResult struct {
	Model      string
	Candidates []PhraseCandidate
}

// Generator proposes candidates: words for a concept, and phrases for a line
// of the phrase list. Claude and Google Translate are the two.
type Generator interface {
	Generate(ctx context.Context, r Request) (Result, error)
	Phrase(ctx context.Context, r PhraseRequest) (PhraseResult, error)
}

// ErrRefused means the model declined. Retrying will not help.
var ErrRefused = errors.New("the model declined the request")

// Lexicon is what this module needs from the lexicon module.
type Lexicon interface {
	Concept(ctx context.Context, id uuid.UUID) (lexicon.Concept, error)
	Add(ctx context.Context, p lexicon.Proposal) (lexicon.Added, error)
	CoreConcepts(ctx context.Context) ([]uuid.UUID, error)
}

var newLexicon = func(tx pgx.Tx) Lexicon { return lexicon.NewStore(tx) }

// Service generates and stores candidates from one provider.
type Service struct {
	db       database.Beginner
	gen      Generator
	sourceID uuid.UUID // the provider's source register row
	actor    string    // audit actor, e.g. "ai:claude"
}

// NewService binds a generator to the source row its candidates cite.
func NewService(db database.Beginner, gen Generator, sourceID uuid.UUID, actor string) *Service {
	return &Service{db: db, gen: gen, sourceID: sourceID, actor: actor}
}

// Stats counts what one call stored: proposals, and the new senses or
// phrases among them. The rest agreed with an item already there.
type Stats struct {
	Proposed, New int
}

// Generate asks the model for one concept in one language and stores the
// answer. The model call happens outside any transaction, since it can take
// minutes.
func (s *Service) Generate(ctx context.Context, conceptID uuid.UUID, language string) (Stats, error) {
	var concept lexicon.Concept
	err := database.InTx(ctx, s.db, s.actor, func(tx pgx.Tx) error {
		var err error
		concept, err = newLexicon(tx).Concept(ctx, conceptID)
		return err
	})
	if err != nil {
		return Stats{}, err
	}

	res, err := s.gen.Generate(ctx, Request{Language: language, Concept: concept})
	if err != nil {
		return Stats{}, err
	}

	stats := Stats{Proposed: len(res.Candidates)}
	origin := items.Machine(s.sourceID, res.Model, "concept:"+conceptID.String())
	err = database.InTx(ctx, s.db, s.actor, func(tx pgx.Tx) error {
		lex := newLexicon(tx)
		for _, c := range res.Candidates {
			added, err := lex.Add(ctx, lexicon.Proposal{
				Language: language, Headword: c.Headword, PartOfSpeech: c.PartOfSpeech,
				NounClass: nounClass(c.NounClass), Plural: c.Plural, GlossEn: concept.GlossEn,
				ConceptID: uuid.NullUUID{UUID: conceptID, Valid: true}, Origin: origin,
			})
			if err != nil {
				return fmt.Errorf("candidate %q: %w", c.Headword, err)
			}
			if added.Sense {
				stats.New++
			}
		}
		return nil
	})
	return stats, err
}

// GeneratePhrase asks the model for one line of the phrase list in one
// language and stores each proposal as a machine phrase (AI-1, PHR-1). It
// returns how many were proposed and how many were new phrases.
func (s *Service) GeneratePhrase(ctx context.Context, lineID uuid.UUID, language string) (Stats, error) {
	var line phrases.Line
	err := database.InTx(ctx, s.db, s.actor, func(tx pgx.Tx) error {
		var err error
		line, err = phrases.NewStore(tx).Line(ctx, lineID)
		return err
	})
	if err != nil {
		return Stats{}, err
	}

	res, err := s.gen.Phrase(ctx, PhraseRequest{Language: language, Pack: line.Pack, SourceEn: line.SourceEn, Notes: line.Notes})
	if err != nil {
		return Stats{}, err
	}

	stats := Stats{Proposed: len(res.Candidates)}
	origin := items.Machine(s.sourceID, res.Model, "phrase:"+line.ID.String())
	err = database.InTx(ctx, s.db, s.actor, func(tx pgx.Tx) error {
		store := phrases.NewStore(tx)
		for _, c := range res.Candidates {
			added, err := store.Propose(ctx, language, line, c.Text, c.Notes, origin)
			if err != nil {
				return fmt.Errorf("phrase candidate %q: %w", c.Text, err)
			}
			if added {
				stats.New++
			}
		}
		return nil
	})
	return stats, err
}

// nounClassPattern is the entries.noun_class check ("1", "1/2").
var nounClassPattern = regexp.MustCompile(`^[0-9]{1,2}(/[0-9]{1,2})?$`)

// nounClass drops a malformed class hint rather than failing the candidate:
// the word is still worth reviewing.
func nounClass(v string) string {
	if v = strings.TrimSpace(v); nounClassPattern.MatchString(v) {
		return v
	}
	return ""
}
