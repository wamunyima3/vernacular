package phrases

import (
	"context"
	"errors"
	"io"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/phrases/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/csvx"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/textnorm"
)

// ImportList loads the English source list of the packs from a CSV with the
// header pack,source_en,notes (P1-API-07b). Each line is linked to the
// concept it expresses, found or created without joining the core list.
// Re-running is idempotent. It returns how many lines are new.
func ImportList(ctx context.Context, db database.Beginner, r io.Reader) (int, error) {
	added := 0
	err := database.InTx(ctx, db, "ingest:phrases", func(tx pgx.Tx) error {
		lex, q := lexicon.NewStore(tx), dbgen.New(tx)
		return csvx.Each(r, []string{"pack", "source_en", "notes"}, func(rec []string) error {
			source := textnorm.NFC(textnorm.Squash(rec[1]))
			concept, err := lex.ConceptFor(ctx, source)
			if err != nil {
				return err
			}
			isNew, err := q.UpsertListLine(ctx, dbgen.UpsertListLineParams{
				Pack: textnorm.Squash(rec[0]), SourceEn: source, Notes: textnorm.Squash(rec[2]), ConceptID: concept,
			})
			if errors.Is(err, pgx.ErrNoRows) {
				return nil // unchanged
			}
			if isNew {
				added++
			}
			return err
		})
	})
	return added, err
}
