// Package contributions lets registered contributors propose words (CON-1).
// Proposals are stored as `sourced`, cite the community contributors source
// and their author, and wait for blind review like any import.
package contributions

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/languages"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

// SourceName is the register row contributions cite (DQ-1).
const SourceName = "Community contributors"

// Input is one proposed word and meaning, with an optional example.
type Input struct {
	Language, Headword, PartOfSpeech, NounClass, GlossEn string
	Example                                              *Example
	// ConceptID names the meaning when it is known, as for review corrections.
	ConceptID uuid.NullUUID
	// Corrects names the sense this proposal would replace (US-18).
	Corrects uuid.NullUUID
}

// Example is a sentence using the word, with its English translation.
type Example struct {
	Text, TranslationEn string
}

// Receipt says what was stored and how far it is trusted.
type Receipt struct {
	EntryID, SenseID uuid.UUID
	ExampleID        *uuid.UUID
	NewSense         bool
	Trust            items.Trust
}

// Lexicon is what contributions need from the lexicon module.
type Lexicon interface {
	Add(ctx context.Context, p lexicon.Proposal) (lexicon.Added, error)
	AddExample(ctx context.Context, language string, senseID uuid.UUID, text, translationEn string, origin items.Origin) (uuid.UUID, error)
}

// newLexicon gives a lexicon bound to the contribution's transaction.
var newLexicon = func(tx pgx.Tx) Lexicon { return lexicon.NewStore(tx) }

// Service records contributions.
type Service struct {
	db database.Beginner
}

// NewService returns a contributions service.
func NewService(db database.Beginner) *Service { return &Service{db: db} }

// Propose stores a contribution from the signed-in contributor.
func (s *Service) Propose(ctx context.Context, in Input) (Receipt, error) {
	author, err := accounts.Require(ctx, func(p accounts.Principal) bool { return p.HasRole(gen.RoleContributor) })
	if err != nil {
		return Receipt{}, err
	}
	var r Receipt
	err = database.InTx(ctx, s.db, author.ID, func(tx pgx.Tx) error {
		if r, err = Add(ctx, tx, author.ID, in); err != nil {
			return err
		}
		r.Trust, err = items.NewStore(tx).Trust(ctx, r.SenseID)
		return err
	})
	return r, err
}

// ExampleReceipt is a contributed example and its trust.
type ExampleReceipt struct {
	ExampleID uuid.UUID
	Trust     items.Trust
}

// ProposeExample adds the signed-in contributor's example to an existing
// meaning and queues it for review (LEX-5).
func (s *Service) ProposeExample(ctx context.Context, senseID uuid.UUID, ex Example) (ExampleReceipt, error) {
	author, err := accounts.Require(ctx, func(p accounts.Principal) bool { return p.HasRole(gen.RoleContributor) })
	if err != nil {
		return ExampleReceipt{}, err
	}
	var r ExampleReceipt
	err = database.InTx(ctx, s.db, author.ID, func(tx pgx.Tx) error {
		its := items.NewStore(tx)
		sense, ok, err := its.Lock(ctx, senseID)
		switch {
		case err != nil:
			return err
		case !ok || sense.Kind != items.KindSense:
			return httpx.ErrNotFound
		}
		if r.ExampleID, err = addExample(ctx, tx, author.ID, sense.Language, senseID, ex); err != nil {
			return err
		}
		r.Trust, err = its.Trust(ctx, r.ExampleID)
		return err
	})
	return r, err
}

// Add stores a contribution by author inside the caller's transaction; the
// caller has checked that author may contribute. Receipt.Trust is left empty.
func Add(ctx context.Context, tx pgx.Tx, author string, in Input) (Receipt, error) {
	if err := languages.NewStore(tx).RequireActive(ctx, in.Language); err != nil {
		return Receipt{}, err
	}
	sourceID, err := sources.Lookup(ctx, tx, SourceName)
	if err != nil {
		return Receipt{}, err
	}
	origin := items.Contributed(sourceID, author)
	lex := newLexicon(tx)
	added, err := lex.Add(ctx, lexicon.Proposal{
		Language: in.Language, Headword: in.Headword, PartOfSpeech: in.PartOfSpeech,
		NounClass: in.NounClass, GlossEn: in.GlossEn, ConceptID: in.ConceptID, Origin: origin,
	})
	if errors.Is(err, lexicon.ErrInvalid) {
		return Receipt{}, httpx.BadRequest(err.Error())
	}
	if err != nil {
		return Receipt{}, err
	}
	r := Receipt{EntryID: added.EntryID, SenseID: added.SenseID, NewSense: added.Sense}
	if in.Corrects.Valid {
		if err := linkCorrection(ctx, tx, in, added); err != nil {
			return Receipt{}, err
		}
	}
	if in.Example != nil {
		id, err := exampleBy(ctx, lex, in.Language, added.SenseID, *in.Example, origin)
		if err != nil {
			return Receipt{}, err
		}
		r.ExampleID = &id
	}
	return r, nil
}

// addExample adds author's example to a meaning in an active language.
func addExample(ctx context.Context, tx pgx.Tx, author, language string, senseID uuid.UUID, ex Example) (uuid.UUID, error) {
	if err := languages.NewStore(tx).RequireActive(ctx, language); err != nil {
		return uuid.Nil, err
	}
	sourceID, err := sources.Lookup(ctx, tx, SourceName)
	if err != nil {
		return uuid.Nil, err
	}
	return exampleBy(ctx, newLexicon(tx), language, senseID, ex, items.Contributed(sourceID, author))
}

// exampleBy stores an example with its origin, a bad one as a 400.
func exampleBy(ctx context.Context, lex Lexicon, language string, senseID uuid.UUID, ex Example, origin items.Origin) (uuid.UUID, error) {
	id, err := lex.AddExample(ctx, language, senseID, ex.Text, ex.TranslationEn, origin)
	if errors.Is(err, lexicon.ErrInvalid) {
		return uuid.Nil, httpx.BadRequest(err.Error())
	}
	return id, err
}

// linkCorrection checks the sense a proposal corrects and links the new
// candidate to it. A correction must say something the dictionary does not.
func linkCorrection(ctx context.Context, tx pgx.Tx, in Input, added lexicon.Added) error {
	its := items.NewStore(tx)
	target, ok, err := its.Lock(ctx, in.Corrects.UUID)
	switch {
	case err != nil:
		return err
	case !ok || target.Kind != items.KindSense || target.Language != in.Language:
		return httpx.BadRequest("corrects must name a meaning in " + in.Language)
	case added.SenseID == in.Corrects.UUID:
		return httpx.BadRequest("That is what the dictionary already says; change the word or meaning")
	case !added.Sense:
		return httpx.BadRequest("That word and meaning are already in the dictionary; report the problem instead")
	}
	return its.SetCorrects(ctx, added.SenseID, in.Corrects)
}
