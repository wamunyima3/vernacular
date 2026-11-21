package detect

import (
	"context"
	"fmt"

	"github.com/wamunyima3/vernacular/apps/api/internal/detect/dbgen"
)

// MinChars is how much text a language needs before the model includes it.
const MinChars = 2000

// Load trains a model from the database (ADR-0010).
func Load(ctx context.Context, db dbgen.DBTX) (*Model, error) {
	rows, err := dbgen.New(db).ListTrainingText(ctx)
	if err != nil {
		return nil, fmt.Errorf("list training text: %w", err)
	}
	samples := map[string][]string{}
	for _, r := range rows {
		samples[r.Language] = append(samples[r.Language], r.Text)
	}
	return Train(samples, MinChars), nil
}
