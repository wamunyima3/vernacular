// Package server composes the module handlers into the generated API and
// mounts it on the shared router. It is the only place modules meet.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	validator "github.com/oapi-codegen/nethttp-middleware"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/audio"
	"github.com/wamunyima3/vernacular/apps/api/internal/contributions"
	"github.com/wamunyima3/vernacular/apps/api/internal/detect"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/health"
	"github.com/wamunyima3/vernacular/apps/api/internal/languages"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/machine"
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
	audioOps     = audio.Handler
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
	*audioOps
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
	Usage         *accounts.Usage // optional: counts calls per key
	Audio         *audio.Service
	Machine       machine.Translator // optional: machine translation (TRN-3)
}

// countUsage records calls per key when deps has a recorder.
func countUsage(u *accounts.Usage) func(http.Handler) http.Handler {
	if u == nil {
		return func(next http.Handler) http.Handler { return next }
	}
	return u.Middleware
}

// New returns the full HTTP handler for the API.
func New(opts Options, deps Deps) http.Handler {
	errs := httpx.Errors{DocsBaseURL: opts.DocsBaseURL, Log: opts.Log}
	r := httpx.NewRouter(httpx.RouterOptions{
		Log: opts.Log, AllowedOrigins: opts.AllowedOrigins, Errors: errs, Metrics: opts.Metrics,
	})
	spec, err := gen.GetSpec()
	if err != nil {
		panic(fmt.Sprintf("embedded OpenAPI spec is invalid: %v", err)) // a build defect, caught by every test
	}
	r.Use(accounts.Middleware(deps.Accounts, errs), countUsage(deps.Usage), rateLimit(opts.RateLimits, errs), httpx.ETag(cacheControl), validateRequests(spec, errs))

	strict := gen.NewStrictHandlerWithOptions(api{
		accountsOps:  accounts.NewHandler(deps.Reviewers, opts.RateLimits.PerKey),
		healthOps:    health.NewHandler(deps.DB, errs),
		languagesOps: languages.NewHandler(deps.Languages),
		sourcesOps:   sources.NewHandler(deps.Sources),
		contribOps:   contributions.NewHandler(deps.Contributions),
		reviewOps:    review.NewHandler(deps.Review),
		lexiconOps:   lexicon.NewHandler(deps.Lexicon),
		phrasesOps:   phrases.NewHandler(deps.Phrases, deps.Machine),
		detectOps:    detect.NewHandler(deps.Detector),
		releasesOps:  releases.NewHandler(deps.Releases),
		audioOps:     audio.NewHandler(deps.Audio),
	}, []gen.StrictMiddlewareFunc{pinRelease(spec, deps.Releases)}, gen.StrictHTTPServerOptions{
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
func validateRequests(spec *openapi3.T, errs httpx.Errors) func(http.Handler) http.Handler {
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

// pinRelease runs the operations whose contract takes `release` on that
// release's frozen rows when the caller names one (ADR-0012). Stores read
// through database.Routed, so they follow the pinned transaction.
func pinRelease(spec *openapi3.T, svc *releases.Service) gen.StrictMiddlewareFunc {
	pinnable := operationsWithParam(spec, "release")
	return func(f gen.StrictHandlerFunc, op string) gen.StrictHandlerFunc {
		if !pinnable[op] {
			return f
		}
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, req any) (resp any, err error) {
			name := r.URL.Query().Get("release")
			if name == "" {
				return f(ctx, w, r, req)
			}
			err = svc.Pin(ctx, name, func(ctx context.Context, _ string) error {
				resp, err = f(ctx, w, r, req)
				return err
			})
			return resp, err
		}
	}
}

// operationsWithParam names the operations whose contract declares the query parameter.
func operationsWithParam(spec *openapi3.T, name string) map[string]bool {
	ops := map[string]bool{}
	for _, item := range spec.Paths.Map() {
		for _, op := range item.Operations() {
			for _, p := range op.Parameters {
				if p.Value.In == "query" && p.Value.Name == name {
					ops[op.OperationID] = true
				}
			}
		}
	}
	return ops
}

// cacheControl lets clients reuse data for five minutes, or a day when it
// is pinned to a numbered release, then revalidate with the ETag; what a signed-in user sees about their own account is
// always revalidated. Private either way, so a shared cache cannot hand data
// fetched with a key to a caller without one.
func cacheControl(r *http.Request) string {
	if p, ok := accounts.FromContext(r.Context()); ok && !p.ByKey() {
		return "private, no-cache"
	}
	if release := r.URL.Query().Get("release"); release != "" && release != "latest" {
		return "private, max-age=86400" // a numbered release's content is frozen
	}
	return "private, max-age=300"
}
