package review

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/contributions"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/textnorm"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

// Correction is a reviewer's own word for a card's meaning (US-22).
type Correction struct {
	Headword, PartOfSpeech, NounClass string
}

// replacement is the candidate a correction created, and the review detail
// that points the corrected candidates at it.
type replacement struct {
	senseID uuid.UUID
	detail  []byte
}

func isCorrection(d Decision) bool { return d.Decision == dbgen.ReviewDecisionCorrect }

// replace adds the reviewer's word as their contribution to the meaning of
// the corrected senses, so only others review it; the corrected senses keep
// their history untouched.
func replace(ctx context.Context, tx pgx.Tx, reviewerID, language string, decisions []Decision, c Correction) (*replacement, error) {
	var ids []uuid.UUID
	for _, d := range decisions {
		if isCorrection(d) {
			ids = append(ids, d.SenseID)
		}
	}
	lex := newLexicon(tx)
	cards, err := lex.SenseCards(ctx, ids)
	if err != nil {
		return nil, err
	}
	// The same word as a candidate, however spelled, should be accepted instead.
	word, _ := lexicon.Standardise(language, textnorm.Headword(c.Headword))
	concept := cards[0].ConceptID
	for _, card := range cards {
		if card.ConceptID != concept || !concept.Valid {
			return nil, httpx.BadRequest("A correction applies to the candidates of one meaning")
		}
		if card.Headword == word {
			return nil, httpx.BadRequest("That word is one of the candidates; accept it instead")
		}
	}
	meaning, err := lex.Concept(ctx, concept.UUID)
	if err != nil {
		return nil, err
	}
	added, err := contributions.Add(ctx, tx, reviewerID, contributions.Input{
		Language: language, Headword: c.Headword, PartOfSpeech: c.PartOfSpeech, NounClass: c.NounClass,
		GlossEn: meaning.GlossEn, ConceptID: concept,
	})
	if err != nil {
		return nil, err
	}
	detail, err := json.Marshal(map[string]uuid.UUID{"sense_id": added.SenseID})
	if err != nil {
		return nil, fmt.Errorf("encode correction: %w", err)
	}
	return &replacement{senseID: added.SenseID, detail: detail}, nil
}
