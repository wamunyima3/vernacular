// Command worker runs background jobs with River. Today that is AI candidate
// generation (P1-AI-01):
//
//	worker                 run jobs until interrupted
//	worker -enqueue bem                   queue Claude jobs for every core concept, then exit
//	worker -enqueue nya -provider google  the same with Google Translate
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

	"cloud.google.com/go/translate"
	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"google.golang.org/api/option"

	"github.com/wamunyima3/vernacular/apps/api/internal/candidates"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

type config struct {
	DatabaseURL    string     `env:"DATABASE_URL,required"`
	LogLevel       slog.Level `env:"LOG_LEVEL" envDefault:"info"`
	StyleGuidesDir string     `env:"STYLE_GUIDES_DIR" envDefault:"../../docs/style-guides"`
	// Google Translate is used only when a key is set.
	GoogleTranslateAPIKey string `env:"GOOGLE_TRANSLATE_API_KEY"`
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
	provider := flag.String("provider", candidates.ProviderClaude, "claude or google")
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
		n, err := candidates.EnqueueCore(ctx, pool, client, *provider, *enqueue)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s: %d jobs queued\n", *provider, *enqueue, n)
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
	providers, err := newProviders(ctx, pool, cfg)
	if err != nil {
		return nil, err
	}
	workers := river.NewWorkers()
	river.AddWorker(workers, candidates.NewWorker(providers, log))
	return river.NewClient(riverpgxv5.New(pool), &river.Config{
		Logger:  log,
		Queues:  map[string]river.QueueConfig{candidates.Queue: {MaxWorkers: cfg.AIConcurrency}},
		Workers: workers,
	})
}

// newProviders builds a candidates service per configured provider, each
// citing its own source register row (DQ-1).
func newProviders(ctx context.Context, pool *pgxpool.Pool, cfg config) (map[string]*candidates.Service, error) {
	gens := map[string]candidates.Generator{}
	claude, err := candidates.NewClaude(anthropic.NewClient(), cfg.StyleGuidesDir)
	if err != nil {
		return nil, err
	}
	gens[candidates.ProviderClaude] = claude
	if cfg.GoogleTranslateAPIKey != "" {
		client, err := translate.NewClient(ctx, option.WithAPIKey(cfg.GoogleTranslateAPIKey))
		if err != nil {
			return nil, fmt.Errorf("google translate client: %w", err)
		}
		gens[candidates.ProviderGoogle] = candidates.NewGoogle(client)
	}

	providers := map[string]*candidates.Service{}
	for name, gen := range gens {
		sourceID, err := sources.Lookup(ctx, pool, candidates.SourceNames[name])
		if err != nil {
			return nil, err
		}
		providers[name] = candidates.NewService(pool, gen, sourceID, "ai:"+name)
	}
	return providers, nil
}
