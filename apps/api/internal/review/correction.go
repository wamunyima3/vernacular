package review

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/contributions"
	"github.com/wamunyima3/vernacular/apps/api/internal/items"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/phrases"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/textnorm"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

// Correction is a reviewer's own word for a card's meaning, phrase for its
// line, or sentence for its examples (US-22).
type Correction struct {
	Headword, PartOfSpeech, NounClass string
	TranslationEn                     string // an example sentence's English
}

// replacement is the candidate a correction created, and the review detail
// that points the corrected candidates at it.
type replacement struct {
	itemID uuid.UUID
	detail []byte
}

// correctionDetail is what a `correct` review stores: the candidate it
// created, a sense, phrase or example (the name predates the others).
type correctionDetail struct {
	ItemID uuid.UUID `json:"sense_id"` // the key predates phrases; stored reviews keep it
}

// correctedTo reads the candidate a `correct` review created; Nil otherwise.
func correctedTo(detail []byte) uuid.UUID {
	var c correctionDetail
	if detail == nil || json.Unmarshal(detail, &c) != nil {
		return uuid.Nil
	}
	return c.ItemID
}

func isCorrection(d Decision) bool { return d.Decision == dbgen.ReviewDecisionCorrect }

// correctedIDs are the items the decisions correct.
func correctedIDs(decisions []Decision) []uuid.UUID {
	var ids []uuid.UUID
	for _, d := range decisions {
		if isCorrection(d) {
			ids = append(ids, d.ItemID)
		}
	}
	return ids
}

// replacePhrase adds the reviewer's phrase as their contribution to the line
// the corrected phrases translate, so only others review it.
func replacePhrase(ctx context.Context, tx pgx.Tx, reviewerID, language string, decisions []Decision, c Correction) (*replacement, error) {
	phr := newPhrases(tx)
	var line phrases.Line
	for _, d := range decisions {
		if !isCorrection(d) {
			continue
		}
		l, err := phr.LineOf(ctx, d.ItemID)
		if err != nil {
			return nil, err
		}
		if line.ID != uuid.Nil && l.ID != line.ID {
			return nil, httpx.BadRequest("A correction applies to the candidates of one phrase")
		}
		line = l
	}
	sourceID, err := sources.Lookup(ctx, tx, contributions.SourceName)
	if err != nil {
		return nil, err
	}
	id, added, err := phr.Propose(ctx, language, line, c.Headword, "", items.Contributed(sourceID, reviewerID))
	switch {
	case errors.Is(err, phrases.ErrPlaceholders):
		return nil, httpx.BadRequest(phrases.PlaceholderMessage(err))
	case err != nil:
		return nil, err
	case id == uuid.Nil:
		return nil, httpx.BadRequest("Write the phrase")
	case !added:
		return nil, httpx.BadRequest("That phrase is one of the candidates; accept it instead")
	}
	return correctedBy(id)
}

// replaceExample adds the reviewer's sentence as their example of the meaning
// the corrected examples illustrate, so only others review it.
func replaceExample(ctx context.Context, tx pgx.Tx, reviewerID, language string, decisions []Decision, c Correction) (*replacement, error) {
	lex := newLexicon(tx)
	cards, err := lex.ExampleCards(ctx, correctedIDs(decisions))
	if err != nil {
		return nil, err
	}
	for _, card := range cards {
		if card.SenseID != cards[0].SenseID {
			return nil, httpx.BadRequest("A correction applies to the examples of one word")
		}
		if card.Text == textnorm.Squash(textnorm.NFC(c.Headword)) {
			return nil, httpx.BadRequest("That sentence is one of the candidates; accept it instead")
		}
	}
	sourceID, err := sources.Lookup(ctx, tx, contributions.SourceName)
	if err != nil {
		return nil, err
	}
	id, err := lex.AddExample(ctx, language, cards[0].SenseID, c.Headword, c.TranslationEn, items.Contributed(sourceID, reviewerID))
	if errors.Is(err, lexicon.ErrInvalid) {
		return nil, httpx.BadRequest("Write the sentence and its English translation")
	}
	if err != nil {
		return nil, err
	}
	return correctedBy(id)
}

// correctedBy is the replacement pointing corrected candidates at id.
func correctedBy(id uuid.UUID) (*replacement, error) {
	detail, err := json.Marshal(correctionDetail{ItemID: id})
	if err != nil {
		return nil, fmt.Errorf("encode correction: %w", err)
	}
	return &replacement{itemID: id, detail: detail}, nil
}

// replace adds the reviewer's word as their contribution to the meaning of
// the corrected senses, so only others review it; the corrected senses keep
// their history untouched.
func replace(ctx context.Context, tx pgx.Tx, reviewerID, language string, decisions []Decision, c Correction) (*replacement, error) {
	lex := newLexicon(tx)
	cards, err := lex.SenseCards(ctx, correctedIDs(decisions))
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
	added, err := contributions.Add(ctx, tx, reviewerID, contributions.Input{
		Language: language, Headword: c.Headword, PartOfSpeech: c.PartOfSpeech, NounClass: c.NounClass,
		GlossEn: cards[0].GlossEn, ConceptID: concept,
	})
	if err != nil {
		return nil, err
	}
	return correctedBy(added.SenseID)
}
