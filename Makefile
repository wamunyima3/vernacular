# Single entry point for Go and JS tasks. Run `make help` for the list.
SHELL := /bin/bash
.DEFAULT_GOAL := help

-include .env
export

DATABASE_URL ?= postgres://vernacular:vernacular@localhost:5432/vernacular?sslmode=disable
# Development object store: a local folder (.data is git-ignored).
BLOB_URL     ?= file://$(CURDIR)/.data/blobs?create_dir=true
API_DIR      := apps/api
# Go dev tools are pinned in apps/api/tools/go.mod, outside the app's own go.mod.
GO_TOOL      := go tool -modfile=tools/go.mod
COMPOSE      := docker compose -f deploy/docker-compose.yml
GOOSE        := $(GO_TOOL) goose -dir db/migrations postgres "$(DATABASE_URL)"
# River (background jobs) owns its tables and migrations.
RIVER        := $(GO_TOOL) river migrate-up --line main --database-url "$(DATABASE_URL)"
# Everything `make gen` writes. Never edited by hand.
GENERATED    := $(API_DIR)/internal/gen $(API_DIR)/internal/*/dbgen packages/api-client/src/schema.d.ts

.PHONY: help
help: ## Show this help
	@grep -hE '^[a-z][a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup: ## Install JS deps and git hooks, download Go modules
	pnpm install
	cd $(API_DIR) && go mod download && go mod download -modfile=tools/go.mod
	pnpm exec lefthook install

.PHONY: up
up: ## Start Postgres and wait until it is healthy
	$(COMPOSE) up -d --wait

.PHONY: down
down: ## Stop local services
	$(COMPOSE) down

.PHONY: migrate
migrate: ## Apply database migrations (ours, then River's)
	cd $(API_DIR) && $(GOOSE) up && $(RIVER)

.PHONY: migrate-down
migrate-down: ## Roll back the last migration
	cd $(API_DIR) && $(GOOSE) down

.PHONY: migrate-roundtrip
migrate-roundtrip: ## Prove every migration applies, rolls back and re-applies
	cd $(API_DIR) && $(GOOSE) up && $(GOOSE) reset && $(GOOSE) up

.PHONY: ingest
ingest: ## Import a registered source, e.g. make ingest FEED=gatitos-en-bem
	cd $(API_DIR) && go run ./cmd/ingest -feed $(FEED)

.PHONY: restem
restem: ## Recompute search stems after changing a stemmer
	cd $(API_DIR) && go run ./cmd/ingest -restem

.PHONY: seed
seed: ## Import our core concept lists and every open source we hold a licence for
	cd $(API_DIR) && go run ./cmd/ingest -all
	for f in data/concepts/*.csv; do (cd $(API_DIR) && go run ./cmd/ingest -concepts ../../$$f) || exit 1; done

.PHONY: admin
admin: ## Give a signed-up user a role, e.g. make admin EMAIL=you@example.org [ROLE=admin]
	cd $(API_DIR) && go run ./cmd/admin -email "$(EMAIL)" -role "$(or $(ROLE),admin)"

.PHONY: gen
gen: ## Regenerate code from the OpenAPI contract and SQL queries
	cd $(API_DIR) && $(GO_TOOL) oapi-codegen -config oapi-codegen.yaml ../../contracts/openapi.yaml
	cd $(API_DIR) && $(GO_TOOL) sqlc generate
	pnpm --filter @vernacular/api-client gen

.PHONY: gen-check
gen-check: gen ## Fail if generated code is out of date
	@stale="$$(git status --porcelain -- $(GENERATED))"; \
	if [ -n "$$stale" ]; then echo "$$stale"; echo "Generated code is stale: run 'make gen' and commit."; exit 1; fi

.PHONY: dev
dev: ## Run the API and the web app together
	$(MAKE) -j2 dev-api dev-web

.PHONY: dev-api
dev-api:
	cd $(API_DIR) && go run ./cmd/api

.PHONY: dev-worker
dev-worker: ## Run background jobs (needs ANTHROPIC_API_KEY)
	cd $(API_DIR) && go run ./cmd/worker

PROVIDER ?= claude
.PHONY: candidates
candidates: ## Queue AI candidate jobs for every core concept: make candidates LANG=bem [PROVIDER=google]
	cd $(API_DIR) && go run ./cmd/worker -enqueue $(LANG) -provider $(PROVIDER)

.PHONY: dev-web
dev-web:
	pnpm --filter @vernacular/web dev

.PHONY: test
test: test-api ## Run Go and web tests
	pnpm test

.PHONY: test-api
test-api:
	cd $(API_DIR) && go test -race ./...

.PHONY: lint
lint: lint-api ## Lint Go and TS, check formatting
	pnpm lint
	pnpm format:check

.PHONY: lint-api
lint-api: # golangci-lint v2 (version pinned in CI); includes go vet
	cd $(API_DIR) && golangci-lint run ./...

.PHONY: typecheck
typecheck: ## Type-check TS (Go is type-checked by build, vet and lint)
	pnpm typecheck

.PHONY: dup
dup: ## Check for copy-pasted code
	pnpm check:dup

.PHONY: check
check: gen-check lint typecheck test dup ## Everything CI runs
	pnpm build
