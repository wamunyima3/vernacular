package candidates

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
)

// Queue is the River queue model calls run on, kept apart so its
// concurrency can match the provider's rate limit.
const Queue = "ai"

// Providers that generate candidates.
const (
	ProviderClaude = "claude"
	ProviderGoogle = "google"
)

// SourceNames maps each provider to the source register row its candidates
// cite (DQ-1).
var SourceNames = map[string]string{
	ProviderClaude: "Claude",
	ProviderGoogle: "Google Translate",
}

// GenerateArgs is one River job: one concept in one language from one provider.
type GenerateArgs struct {
	ConceptID uuid.UUID `json:"concept_id"`
	Language  string    `json:"language"`
	Provider  string    `json:"provider"`
}

// Kind names the job in River.
func (GenerateArgs) Kind() string { return "generate_candidates" }

// InsertOpts makes the job unique per concept and language, so enqueuing
// the whole list again skips the concepts already done.
func (GenerateArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{Queue: Queue, MaxAttempts: 5, UniqueOpts: river.UniqueOpts{ByArgs: true}}
}

// Worker runs GenerateArgs jobs with the service for each provider.
type Worker struct {
	river.WorkerDefaults[GenerateArgs]
	providers map[string]*Service
	log       *slog.Logger
}

// NewWorker wraps one service per provider as a River worker.
func NewWorker(providers map[string]*Service, log *slog.Logger) *Worker {
	return &Worker{providers: providers, log: log}
}

// Timeout allows for a long model turn.
func (*Worker) Timeout(*river.Job[GenerateArgs]) time.Duration { return 10 * time.Minute }

// Work generates and stores candidates. A refusal is final, so the job is
// cancelled rather than retried.
func (w *Worker) Work(ctx context.Context, job *river.Job[GenerateArgs]) error {
	svc, ok := w.providers[job.Args.Provider]
	if !ok {
		return river.JobCancel(fmt.Errorf("provider %q is not configured", job.Args.Provider))
	}
	stats, err := svc.Generate(ctx, job.Args.ConceptID, job.Args.Language)
	if errors.Is(err, ErrRefused) {
		return river.JobCancel(err)
	}
	if err != nil {
		return err
	}
	w.log.InfoContext(ctx, "candidates stored", "provider", job.Args.Provider, "concept_id", job.Args.ConceptID, "language", job.Args.Language,
		"proposed", stats.Proposed, "new_senses", stats.Senses, "agreed", stats.Proposed-stats.Senses)
	return nil
}

// EnqueueCore adds one job per core concept for a language and provider.
// Concepts that already have a job are skipped (InsertOpts), so it can be
// re-run.
func EnqueueCore(ctx context.Context, db database.Beginner, client *river.Client[pgx.Tx], provider, language string) (int, error) {
	inserted := 0
	err := database.InTx(ctx, db, "ai:enqueue", func(tx pgx.Tx) error {
		ids, err := newLexicon(tx).CoreConcepts(ctx)
		if err != nil {
			return err
		}
		params := make([]river.InsertManyParams, len(ids))
		for i, id := range ids {
			params[i] = river.InsertManyParams{Args: GenerateArgs{ConceptID: id, Language: language, Provider: provider}}
		}
		res, err := client.InsertManyTx(ctx, tx, params)
		if err != nil {
			return fmt.Errorf("enqueue: %w", err)
		}
		for _, r := range res {
			if !r.UniqueSkippedAsDuplicate {
				inserted++
			}
		}
		return nil
	})
	return inserted, err
}
