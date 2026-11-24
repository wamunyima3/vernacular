// Package config loads process configuration from environment variables once at startup.
package config

import (
	"log/slog"
	"net/netip"

	"github.com/caarlos0/env/v11"
)

// Config is the API process configuration. Every field is documented in .env.example.
type Config struct {
	Env                string     `env:"APP_ENV"              envDefault:"development"`
	HTTPAddr           string     `env:"HTTP_ADDR"            envDefault:":8080"`
	MetricsAddr        string     `env:"METRICS_ADDR"         envDefault:":9090"` // Prometheus; never public. Empty turns it off.
	LogLevel           slog.Level `env:"LOG_LEVEL"            envDefault:"info"`
	DatabaseURL        string     `env:"DATABASE_URL,required"`
	CORSAllowedOrigins []string   `env:"CORS_ALLOWED_ORIGINS" envDefault:"http://localhost:3000" envSeparator:","`
	DocsBaseURL        string     `env:"DOCS_BASE_URL"        envDefault:"http://localhost:3000/docs/errors"`
	BlobURL            string     `env:"BLOB_URL,required"` // see internal/platform/blob
	// Requests per minute (DEV-3); 0 turns a limit off.
	RateLimitPerKey int `env:"RATE_LIMIT_PER_KEY" envDefault:"120"`
	RateLimitPerIP  int `env:"RATE_LIMIT_PER_IP"  envDefault:"30"`
	// Proxies whose X-Forwarded-For is believed, as CIDRs (NFR-5).
	TrustedProxies []netip.Prefix `env:"TRUSTED_PROXIES" envSeparator:","`
}

// Load parses the environment into a Config.
func Load() (Config, error) {
	return env.ParseAs[Config]()
}
