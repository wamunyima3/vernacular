package sources

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
)

// ImportConcepts adds the project's own concept lists (data/concepts/*.csv,
// columns domain,gloss_en) to the core list. These are English concepts only;
// words for them arrive through AI candidates and review. It returns how many
// concepts were new.
func ImportConcepts(ctx context.Context, db database.Beginner, r io.Reader) (int, error) {
	rows := csv.NewReader(r)
	header, err := rows.Read()
	if err != nil {
		return 0, fmt.Errorf("read header: %w", err)
	}
	if !slices.Equal(header, []string{"domain", "gloss_en"}) {
		return 0, fmt.Errorf("header = %v, want [domain gloss_en]", header)
	}

	added := 0
	err = database.InTx(ctx, db, "ingest:concepts", func(tx pgx.Tx) error {
		lex := newLexicon(tx)
		for {
			rec, err := rows.Read()
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return err
			}
			isNew, err := lex.AddCoreConcept(ctx, rec[1], rec[0])
			if err != nil {
				line, _ := rows.FieldPos(0)
				return fmt.Errorf("line %d: %w", line, err)
			}
			if isNew {
				added++
			}
		}
	})
	return added, err
}
