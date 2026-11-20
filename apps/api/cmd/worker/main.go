// Command worker runs background jobs with River. Today that is AI candidate
// generation (P1-AI-01):
//
//	worker                 run jobs until interrupted
//	worker -enqueue bem    queue one job per core concept for a language, then exit
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"

	"github.com/wamunyima3/vernacular/apps/api/internal/candidates"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

type config struct {
	DatabaseURL    string     `env:"DATABASE_URL,required"`
	LogLevel       slog.Level `env:"LOG_LEVEL" envDefault:"info"`
	StyleGuidesDir string     `env:"STYLE_GUIDES_DIR" envDefault:"../../docs/style-guides"`
	// Parallel model calls. Keep within the API rate limit.
	AIConcurrency int `env:"AI_CONCURRENCY" envDefault:"4"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx); err != nil {
		slog.Error("worker failed", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	enqueue := flag.String("enqueue", "", "queue candidate jobs for every core concept in this language (bem or nya), then exit")
	flag.Parse()
	cfg, err := env.ParseAs[config]()
	if err != nil {
		return err
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if *enqueue != "" {
		client, err := river.NewClient(riverpgxv5.New(pool), &river.Config{})
		if err != nil {
			return err
		}
		n, err := candidates.EnqueueCore(ctx, pool, client, *enqueue)
		if err != nil {
			return err
		}
		fmt.Printf("%s: %d jobs queued\n", *enqueue, n)
		return nil
	}

	client, err := newClient(ctx, pool, cfg, log)
	if err != nil {
		return err
	}
	// Start outlives ctx: an interrupt triggers the soft Stop below, which lets
	// running jobs finish instead of cancelling them.
	if err := client.Start(context.WithoutCancel(ctx)); err != nil {
		return fmt.Errorf("start river: %w", err)
	}
	log.Info("worker started", "queue", candidates.Queue, "concurrency", cfg.AIConcurrency)
	<-ctx.Done()
	if err := client.Stop(context.WithoutCancel(ctx)); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("stop river: %w", err)
	}
	return nil
}

func newClient(ctx context.Context, pool *pgxpool.Pool, cfg config, log *slog.Logger) (*river.Client[pgx.Tx], error) {
	sourceID, err := sources.Lookup(ctx, pool, "Claude")
	if err != nil {
		return nil, err
	}
	claude, err := candidates.NewClaude(anthropic.NewClient(), cfg.StyleGuidesDir)
	if err != nil {
		return nil, err
	}
	workers := river.NewWorkers()
	river.AddWorker(workers, candidates.NewWorker(candidates.NewService(pool, claude, sourceID, "ai:claude"), log))
	return river.NewClient(riverpgxv5.New(pool), &river.Config{
		Logger:  log,
		Queues:  map[string]river.QueueConfig{candidates.Queue: {MaxWorkers: cfg.AIConcurrency}},
		Workers: workers,
	})
}
