# Single entry point for Go and JS tasks. Run `make help` for the list.
SHELL := /bin/bash
.DEFAULT_GOAL := help

-include .env
export

DATABASE_URL ?= postgres://vernacular:vernacular@localhost:5432/vernacular?sslmode=disable
API_DIR      := apps/api
# Go dev tools are pinned in apps/api/tools/go.mod, outside the app's own go.mod.
GO_TOOL      := go tool -modfile=tools/go.mod
COMPOSE      := docker compose -f deploy/docker-compose.yml
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
migrate: ## Apply database migrations
	cd $(API_DIR) && $(GO_TOOL) goose -dir db/migrations postgres "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down: ## Roll back the last migration
	cd $(API_DIR) && $(GO_TOOL) goose -dir db/migrations postgres "$(DATABASE_URL)" down

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

.PHONY: dev-web
dev-web:
	pnpm --filter @vernacular/web dev

.PHONY: test
test: ## Run Go and web tests
	cd $(API_DIR) && go test -race ./...
	pnpm test

.PHONY: lint
lint: ## Lint Go and TS, check formatting
	cd $(API_DIR) && golangci-lint run ./...
	pnpm lint
	pnpm format:check

.PHONY: typecheck
typecheck: ## Type-check Go and TS
	cd $(API_DIR) && go vet ./...
	pnpm typecheck

.PHONY: dup
dup: ## Check for copy-pasted code
	pnpm check:dup

.PHONY: check
check: gen-check lint typecheck test dup ## Everything CI runs
	pnpm build
