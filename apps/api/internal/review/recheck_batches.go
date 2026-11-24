package review

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/review/dbgen"
)

// Batch names a batch in a language: the source and model its items came from.
type Batch struct {
	Language string
	SourceID uuid.UUID
	Model    string
}

// Fails reports whether a batch's re-checks fail too often to keep it:
// more than 2% of those checked (DQ-6).
func Fails(checked, failed int) bool { return failed*100 > checked*2 }

// Batches is a quarter's re-checks per batch, worst first (admins).
func (s *Service) Batches(ctx context.Context, quarter string) ([]dbgen.RecheckBatchesRow, error) {
	me, err := accounts.RequireAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var rows []dbgen.RecheckBatchesRow
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		rows, err = dbgen.New(tx).RecheckBatches(ctx, quarter)
		return err
	})
	return rows, err
}

// SendBackBatch sends every verified item of a batch failing its quarter's
// re-checks back to the review queue, as a send-back resolution does, and
// returns how many (admins). A batch is sent back once per quarter.
func (s *Service) SendBackBatch(ctx context.Context, quarter string, b Batch) (int, error) {
	me, err := accounts.RequireAdmin(ctx)
	if err != nil {
		return 0, err
	}
	var sent int
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		q, store, its := dbgen.New(tx), NewStore(tx), newItems(tx)
		r, err := q.RecheckBatch(ctx, dbgen.RecheckBatchParams{Quarter: quarter, Language: b.Language, SourceID: b.SourceID, Model: b.Model})
		if err != nil {
			return fmt.Errorf("re-checks of a batch: %w", err)
		}
		if !Fails(int(r.Checked), int(r.Failed)) {
			return httpx.Conflict(fmt.Sprintf("%d of %d re-checks failed, not more than 2%%; the batch stays", r.Failed, r.Checked))
		}
		key := dbgen.SendBackBatchParams{Quarter: quarter, Language: b.Language, SourceID: b.SourceID, Model: b.Model, SentBackBy: me.ID}
		if n, err := q.SendBackBatch(ctx, key); err != nil {
			return fmt.Errorf("record batch send-back: %w", err)
		} else if n == 0 {
			return httpx.Conflict("This batch was already sent back for " + quarter)
		}
		ids, err := q.BatchItems(ctx, dbgen.BatchItemsParams{Language: b.Language, SourceID: b.SourceID, Model: b.Model})
		if err != nil {
			return fmt.Errorf("items of a batch: %w", err)
		}
		note := "Batch failed its " + quarter + " re-check"
		for _, id := range ids {
			item, _, err := its.Lock(ctx, id)
			if err != nil {
				return err
			}
			if _, err := sendBack(ctx, store, its, id, item.Origin, me.ID, note); err != nil {
				return err
			}
		}
		sent = len(ids)
		return q.SetBatchItems(ctx, dbgen.SetBatchItemsParams{
			Quarter: quarter, Language: b.Language, SourceID: b.SourceID, Model: b.Model, Items: int32(sent),
		})
	})
	return sent, err
}
