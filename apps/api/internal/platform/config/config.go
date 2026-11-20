// Package config loads process configuration from environment variables once at startup.
package config

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
)

// Config is the API process configuration. Every field is documented in .env.example.
type Config struct {
	Env                string     `env:"APP_ENV"              envDefault:"development"`
	HTTPAddr           string     `env:"HTTP_ADDR"            envDefault:":8080"`
	LogLevel           slog.Level `env:"LOG_LEVEL"            envDefault:"info"`
	DatabaseURL        string     `env:"DATABASE_URL,required"`
	CORSAllowedOrigins []string   `env:"CORS_ALLOWED_ORIGINS" envDefault:"http://localhost:3000" envSeparator:","`
	DocsBaseURL        string     `env:"DOCS_BASE_URL"        envDefault:"http://localhost:3000/docs/errors"`
}

// Load parses the environment into a Config.
func Load() (Config, error) {
	return env.ParseAs[Config]()
}
