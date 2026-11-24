// Command api serves the public Vernacular HTTP API.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/audio"
	"github.com/wamunyima3/vernacular/apps/api/internal/contributions"
	"github.com/wamunyima3/vernacular/apps/api/internal/detect"
	"github.com/wamunyima3/vernacular/apps/api/internal/languages"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/machine"
	"github.com/wamunyima3/vernacular/apps/api/internal/phrases"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/blob"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/config"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/releases"
	"github.com/wamunyima3/vernacular/apps/api/internal/review"
	"github.com/wamunyima3/vernacular/apps/api/internal/server"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

func main() {
	if err := run(); err != nil {
		slog.Error("api stopped", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	// Read stores follow a pinned release's transaction (ADR-0012).
	reads := database.Routed{Default: pool}

	bucket, err := blob.Open(ctx, cfg.BlobURL)
	if err != nil {
		return err
	}
	defer func() { _ = bucket.Close() }()

	lex := lexicon.NewStore(reads) // read-only here: search and counts
	detector, err := detect.Load(ctx, pool)
	if err != nil {
		return err
	}
	log.Info("language detector trained", "languages", detector.Languages())
	usage := accounts.NewUsage(pool)
	go usage.Run(ctx, 30*time.Second, log)
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	handler := server.New(
		server.Options{
			Log: log, AllowedOrigins: cfg.CORSAllowedOrigins, DocsBaseURL: cfg.DocsBaseURL,
			Metrics: httpx.NewMetrics(reg),
			RateLimits: server.RateLimits{
				PerKey: cfg.RateLimitPerKey, PerIP: cfg.RateLimitPerIP,
				IP: httpx.ClientIP{Trusted: cfg.TrustedProxies},
			},
		},
		server.Deps{
			DB:            pool,
			Accounts:      accounts.NewStore(pool),
			Reviewers:     accounts.NewService(pool),
			Languages:     languages.NewService(languages.NewStore(pool), lex),
			Sources:       sources.NewRegister(pool, bucket),
			Contributions: contributions.NewService(pool),
			Review:        review.NewService(pool),
			Lexicon:       lex,
			Phrases:       phrases.NewStore(reads),
			Detector:      detector,
			Releases:      releases.NewService(pool),
			Usage:         usage,
			Audio:         audio.NewService(pool, bucket, cfg.AudioEnabled),
			Machine:       machine.New(cfg.TranslationServiceURL),
		},
	)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 2)
	go func() {
		log.Info("listening", "addr", cfg.HTTPAddr, "env", cfg.Env)
		errCh <- srv.ListenAndServe()
	}()
	metrics := metricsServer(cfg.MetricsAddr, reg)
	if metrics != nil {
		go func() {
			log.Info("metrics listening", "addr", cfg.MetricsAddr)
			errCh <- metrics.ListenAndServe()
		}()
	}

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		log.Info("shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if metrics != nil {
		_ = metrics.Shutdown(shutdownCtx)
	}
	err = srv.Shutdown(shutdownCtx)
	if ferr := usage.Flush(shutdownCtx); ferr != nil { // after the last call, before the pool closes
		log.Error("key usage", "err", ferr)
	}
	return err
}

// metricsServer serves Prometheus metrics on their own port, or nil when addr is empty.
func metricsServer(addr string, reg *prometheus.Registry) *http.Server {
	if addr == "" {
		return nil
	}
	mux := http.NewServeMux()
	mux.Handle("GET /metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	return &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
}
