package sources

import (
	"context"
	"io"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/platform/csvx"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
)

// ImportConcepts adds the project's own concept lists (data/concepts/*.csv,
// columns domain,gloss_en) to the core list. These are English concepts only;
// words for them arrive through AI candidates and review. It returns how many
// concepts were new.
func ImportConcepts(ctx context.Context, db database.Beginner, r io.Reader) (int, error) {
	added := 0
	err := database.InTx(ctx, db, "ingest:concepts", func(tx pgx.Tx) error {
		lex := newLexicon(tx)
		return csvx.Each(r, []string{"domain", "gloss_en"}, func(rec []string) error {
			isNew, err := lex.AddCoreConcept(ctx, rec[1], rec[0])
			if isNew {
				added++
			}
			return err
		})
	})
	return added, err
}
