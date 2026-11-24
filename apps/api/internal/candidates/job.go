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

	"github.com/wamunyima3/vernacular/apps/api/internal/phrases"
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

// insertOpts makes a job unique by its arguments, so enqueuing a whole list
// again skips what is already queued or done.
var insertOpts = river.InsertOpts{Queue: Queue, MaxAttempts: 5, UniqueOpts: river.UniqueOpts{ByArgs: true}}

// GenerateArgs is one River job: one concept in one language from one provider.
type GenerateArgs struct {
	ConceptID uuid.UUID `json:"concept_id"`
	Language  string    `json:"language"`
	Provider  string    `json:"provider"`
}

// Kind names the job in River.
func (GenerateArgs) Kind() string { return "generate_candidates" }

// InsertOpts makes the job unique per concept, language and provider.
func (GenerateArgs) InsertOpts() river.InsertOpts { return insertOpts }

func (a GenerateArgs) provider() string { return a.Provider }

func (a GenerateArgs) run(ctx context.Context, svc *Service) (Stats, []any, error) {
	stats, err := svc.Generate(ctx, a.ConceptID, a.Language)
	return stats, []any{"concept_id", a.ConceptID, "language", a.Language}, err
}

// PhraseArgs is one River job: one line of the phrase list in one language
// from one provider.
type PhraseArgs struct {
	LineID   uuid.UUID `json:"line_id"`
	Language string    `json:"language"`
	Provider string    `json:"provider"`
}

// Kind names the job in River.
func (PhraseArgs) Kind() string { return "generate_phrase_candidates" }

// InsertOpts makes the job unique per line, language and provider.
func (PhraseArgs) InsertOpts() river.InsertOpts { return insertOpts }

func (a PhraseArgs) provider() string { return a.Provider }

func (a PhraseArgs) run(ctx context.Context, svc *Service) (Stats, []any, error) {
	stats, err := svc.GeneratePhrase(ctx, a.LineID, a.Language)
	return stats, []any{"line_id", a.LineID, "language", a.Language}, err
}

// job is what a Worker runs: it names its provider and calls the service.
type job interface {
	river.JobArgs
	provider() string
	run(ctx context.Context, svc *Service) (Stats, []any, error)
}

// Worker runs one kind of candidate job with the service for each provider.
type Worker[A job] struct {
	river.WorkerDefaults[A]
	providers map[string]*Service
	log       *slog.Logger
}

// NewWorker runs word candidate jobs (GenerateArgs).
func NewWorker(providers map[string]*Service, log *slog.Logger) *Worker[GenerateArgs] {
	return &Worker[GenerateArgs]{providers: providers, log: log}
}

// NewPhraseWorker runs phrase candidate jobs (PhraseArgs).
func NewPhraseWorker(providers map[string]*Service, log *slog.Logger) *Worker[PhraseArgs] {
	return &Worker[PhraseArgs]{providers: providers, log: log}
}

// Timeout allows for a long model turn.
func (*Worker[A]) Timeout(*river.Job[A]) time.Duration { return 10 * time.Minute }

// Work generates and stores candidates. A refusal is final, so the job is
// cancelled rather than retried.
func (w *Worker[A]) Work(ctx context.Context, j *river.Job[A]) error {
	svc, ok := w.providers[j.Args.provider()]
	if !ok {
		return river.JobCancel(fmt.Errorf("provider %q is not configured", j.Args.provider()))
	}
	stats, attrs, err := j.Args.run(ctx, svc)
	if errors.Is(err, ErrRefused) {
		return river.JobCancel(err)
	}
	if err != nil {
		return err
	}
	w.log.InfoContext(ctx, "candidates stored", append(attrs, "kind", j.Args.Kind(), "provider", j.Args.provider(),
		"proposed", stats.Proposed, "new", stats.New, "agreed", stats.Proposed-stats.New)...)
	return nil
}

// EnqueueCore adds one job per core concept for a language and provider.
// Concepts that already have a job are skipped (InsertOpts), so it can be
// re-run.
func EnqueueCore(ctx context.Context, db database.Beginner, client *river.Client[pgx.Tx], provider, language string) (int, error) {
	return enqueue(ctx, db, client,
		func(tx pgx.Tx) ([]uuid.UUID, error) { return newLexicon(tx).CoreConcepts(ctx) },
		func(id uuid.UUID) river.JobArgs {
			return GenerateArgs{ConceptID: id, Language: language, Provider: provider}
		})
}

// EnqueuePhrases adds one job per line of the phrase list for a language and
// provider, skipping lines that already have one.
func EnqueuePhrases(ctx context.Context, db database.Beginner, client *river.Client[pgx.Tx], provider, language string) (int, error) {
	return enqueue(ctx, db, client,
		func(tx pgx.Tx) ([]uuid.UUID, error) { return phrases.NewStore(tx).LineIDs(ctx) },
		func(id uuid.UUID) river.JobArgs {
			return PhraseArgs{LineID: id, Language: language, Provider: provider}
		})
}

// enqueue inserts one job per listed ID and returns how many were new.
func enqueue(ctx context.Context, db database.Beginner, client *river.Client[pgx.Tx],
	list func(pgx.Tx) ([]uuid.UUID, error), args func(uuid.UUID) river.JobArgs,
) (int, error) {
	inserted := 0
	err := database.InTx(ctx, db, "ai:enqueue", func(tx pgx.Tx) error {
		ids, err := list(tx)
		if err != nil {
			return err
		}
		params := make([]river.InsertManyParams, len(ids))
		for i, id := range ids {
			params[i] = river.InsertManyParams{Args: args(id)}
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
