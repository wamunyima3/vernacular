// Package contributions lets registered contributors propose words (CON-1).
// Proposals are stored as `sourced`, cite the community contributors source
// and their author, and wait for blind review like any import.
package contributions

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/languages"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

// SourceName is the register row contributions cite (DQ-1).
const SourceName = "Community contributors"

// Input is one proposed word and meaning, with an optional example.
type Input struct {
	Language, Headword, PartOfSpeech, NounClass, GlossEn string
	Example                                              *Example
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
		if err := languages.NewStore(tx).RequireActive(ctx, in.Language); err != nil {
			return err
		}
		sourceID, err := sources.Lookup(ctx, tx, SourceName)
		if err != nil {
			return err
		}
		origin := items.Contributed(sourceID, author.ID)
		lex := newLexicon(tx)
		added, err := lex.Add(ctx, lexicon.Proposal{
			Language: in.Language, Headword: in.Headword, PartOfSpeech: in.PartOfSpeech,
			NounClass: in.NounClass, GlossEn: in.GlossEn, Origin: origin,
		})
		if err != nil {
			return err
		}
		r = Receipt{EntryID: added.EntryID, SenseID: added.SenseID, NewSense: added.Sense}
		if in.Example != nil {
			id, err := lex.AddExample(ctx, in.Language, added.SenseID, in.Example.Text, in.Example.TranslationEn, origin)
			if err != nil {
				return err
			}
			r.ExampleID = &id
		}
		r.Trust, err = items.NewStore(tx).Trust(ctx, added.SenseID)
		return err
	})
	return r, err
}
