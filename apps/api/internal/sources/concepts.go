package sources

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
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
		lex := lexicon.NewStore(tx)
		for {
			rec, err := rows.Read()
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return err
			}
			domain, gloss := headword(rec[0]), gloss(rec[1])
			if domain == "" || gloss == "" {
				line, _ := rows.FieldPos(0)
				return fmt.Errorf("line %d: domain and gloss_en are required", line)
			}
			isNew, err := lex.AddCoreConcept(ctx, gloss, domain)
			if err != nil {
				return err
			}
			if isNew {
				added++
			}
		}
	})
	return added, err
}
