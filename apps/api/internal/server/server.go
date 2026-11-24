// Package server composes the module handlers into the generated API and
// mounts it on the shared router. It is the only place modules meet.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	validator "github.com/oapi-codegen/nethttp-middleware"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/contributions"
	"github.com/wamunyima3/vernacular/apps/api/internal/detect"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/health"
	"github.com/wamunyima3/vernacular/apps/api/internal/languages"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/phrases"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/releases"
	"github.com/wamunyima3/vernacular/apps/api/internal/review"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

// Every module names its handler type Handler, so aliases give each embedded field a unique name.
type (
	accountsOps  = accounts.Handler
	healthOps    = health.Handler
	languagesOps = languages.Handler
	sourcesOps   = sources.Handler
	contribOps   = contributions.Handler
	reviewOps    = review.Handler
	lexiconOps   = lexicon.Handler
	phrasesOps   = phrases.Handler
	detectOps    = detect.Handler
	releasesOps  = releases.Handler
)

// api satisfies gen.StrictServerInterface by embedding one handler per module.
// Adding an operation to the contract fails compilation here until a module implements it.
type api struct {
	*accountsOps
	*healthOps
	*languagesOps
	*sourcesOps
	*contribOps
	*reviewOps
	*lexiconOps
	*phrasesOps
	*detectOps
	*releasesOps
}

var _ gen.StrictServerInterface = api{}

// Options configures the HTTP layer.
type Options struct {
	Log            *slog.Logger
	AllowedOrigins []string
	RateLimits     RateLimits
	DocsBaseURL    string
	Metrics        *httpx.Metrics // optional; see httpx.NewMetrics
}

// Deps are the module dependencies the server wires together.
type Deps struct {
	DB            health.Pinger
	Accounts      accounts.Authenticator
	Reviewers     *accounts.Service
	Languages     languages.Lister
	Sources       *sources.Register
	Contributions *contributions.Service
	Review        *review.Service
	Lexicon       lexicon.Reader
	Phrases       phrases.Reader
	Detector      *detect.Model
	Releases      *releases.Service
}

// New returns the full HTTP handler for the API.
func New(opts Options, deps Deps) http.Handler {
	errs := httpx.Errors{DocsBaseURL: opts.DocsBaseURL, Log: opts.Log}
	r := httpx.NewRouter(httpx.RouterOptions{
		Log: opts.Log, AllowedOrigins: opts.AllowedOrigins, Errors: errs, Metrics: opts.Metrics,
	})
	r.Use(accounts.Middleware(deps.Accounts, errs), rateLimit(opts.RateLimits, errs), httpx.ETag(cacheControl), validateRequests(errs))

	strict := gen.NewStrictHandlerWithOptions(api{
		accountsOps:  accounts.NewHandler(deps.Reviewers),
		healthOps:    health.NewHandler(deps.DB, errs),
		languagesOps: languages.NewHandler(deps.Languages),
		sourcesOps:   sources.NewHandler(deps.Sources),
		contribOps:   contributions.NewHandler(deps.Contributions),
		reviewOps:    review.NewHandler(deps.Review),
		lexiconOps:   lexicon.NewHandler(deps.Lexicon),
		phrasesOps:   phrases.NewHandler(deps.Phrases),
		detectOps:    detect.NewHandler(deps.Detector),
		releasesOps:  releases.NewHandler(deps.Releases),
	}, nil, gen.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  errs.RequestError,
		ResponseErrorHandlerFunc: errs.ResponseError,
	})
	return gen.HandlerWithOptions(strict, gen.ChiServerOptions{
		BaseRouter:       r,
		ErrorHandlerFunc: errs.RequestError,
	})
}

// validateRequests checks every request against the contract before any
// handler runs, so handlers can trust required fields, enums and formats.
// Authentication is left to accounts.Middleware.
func validateRequests(errs httpx.Errors) func(http.Handler) http.Handler {
	spec, err := gen.GetSpec()
	if err != nil {
		panic(fmt.Sprintf("embedded OpenAPI spec is invalid: %v", err)) // a build defect, caught by every test
	}
	spec.Servers = nil // match any host
	codes := map[int]gen.ErrorCode{
		http.StatusUnauthorized:          gen.ErrorCodeUnauthorized,
		http.StatusNotFound:              gen.ErrorCodeNotFound,
		http.StatusMethodNotAllowed:      gen.ErrorCodeMethodNotAllowed,
		http.StatusRequestEntityTooLarge: gen.ErrorCodePayloadTooLarge,
	}
	return validator.OapiRequestValidatorWithOptions(spec, &validator.Options{
		Options: openapi3filter.Options{AuthenticationFunc: authenticate},
		ErrorHandler: func(w http.ResponseWriter, message string, status int) {
			if status == http.StatusUnauthorized {
				message = "Send an API key in the X-API-Key header, or sign in"
			}
			code, ok := codes[status]
			if !ok {
				code, status = gen.ErrorCodeBadRequest, http.StatusBadRequest
			}
			errs.Write(w, status, code, message)
		},
	})
}

// authenticate checks a security scheme of the contract against the caller
// accounts.Middleware found: apiKey needs a key, session a signed-in user.
func authenticate(_ context.Context, in *openapi3filter.AuthenticationInput) error {
	p, ok := accounts.FromContext(in.RequestValidationInput.Request.Context())
	if ok && p.ByKey() == (in.SecuritySchemeName == "apiKey") {
		return nil
	}
	return errNoCredentials
}

var errNoCredentials = errors.New("no credentials for this scheme")

// cacheControl lets clients reuse data for five minutes, then revalidate
// with the ETag; what a signed-in user sees about their own account is
// always revalidated. Private either way, so a shared cache cannot hand data
// fetched with a key to a caller without one.
func cacheControl(r *http.Request) string {
	if p, ok := accounts.FromContext(r.Context()); ok && !p.ByKey() {
		return "private, no-cache"
	}
	return "private, max-age=300"
}
