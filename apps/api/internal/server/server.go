// Package server composes the module handlers into the generated API and
// mounts it on the shared router. It is the only place modules meet.
package server

import (
	"log/slog"
	"net/http"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/health"
	"github.com/wamunyima3/vernacular/apps/api/internal/languages"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// Every module names its handler type Handler, so aliases give each embedded field a unique name.
type (
	accountsOps  = accounts.Handler
	healthOps    = health.Handler
	languagesOps = languages.Handler
)

// api satisfies gen.StrictServerInterface by embedding one handler per module.
// Adding an operation to the contract fails compilation here until a module implements it.
type api struct {
	*accountsOps
	*healthOps
	*languagesOps
}

var _ gen.StrictServerInterface = api{}

// Options configures the HTTP layer.
type Options struct {
	Log            *slog.Logger
	AllowedOrigins []string
	DocsBaseURL    string
}

// Deps are the module dependencies the server wires together.
type Deps struct {
	DB        health.Pinger
	Accounts  accounts.Authenticator
	Languages languages.Lister
}

// New returns the full HTTP handler for the API.
func New(opts Options, deps Deps) http.Handler {
	errs := httpx.Errors{DocsBaseURL: opts.DocsBaseURL, Log: opts.Log}
	r := httpx.NewRouter(httpx.RouterOptions{Log: opts.Log, AllowedOrigins: opts.AllowedOrigins, Errors: errs})
	r.Use(accounts.Middleware(deps.Accounts, errs))

	strict := gen.NewStrictHandlerWithOptions(api{
		accountsOps:  accounts.NewHandler(),
		healthOps:    health.NewHandler(deps.DB, errs),
		languagesOps: languages.NewHandler(deps.Languages),
	}, nil, gen.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  errs.RequestError,
		ResponseErrorHandlerFunc: errs.ResponseError,
	})
	return gen.HandlerWithOptions(strict, gen.ChiServerOptions{
		BaseRouter:       r,
		ErrorHandlerFunc: errs.RequestError,
	})
}
