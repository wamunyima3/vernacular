# AGENTS.md

Rules for everyone who changes this repo: developers and AI coding agents alike. If a rule here conflicts with a habit, the rule wins. If a rule is wrong, change it in a commit of its own.

Read first: [tasks.md](tasks.md) → [docs/architecture.md](docs/architecture.md) → [docs/requirements.md](docs/requirements.md) for the area you touch.

## 1. Product rules that code must never break

1. **Every returned item carries trust fields:** `status`, `sources` (with licence), `reviews` and `release`. Reuse the `Trust` schema; never write a new one.
2. **Only two independent human accepts make an item `verified`.** The single promotion function lives in the `review` module. Nothing else writes `status = 'verified'`.
3. **`machine` items are returned only when the caller asks, and they are always labelled.** `disputed` items are never returned.
4. **No source is ingested without a source register row** (DQ-1). Never scrape apps or sites that have not given permission.
5. **The test set is never used as training data or as a source of entries** (DQ-10).
6. **Personal data** (contributors, speakers, audio) needs recorded consent and must be deletable (NFR-6).

## 2. How to pick and run a task

1. Take the first unchecked task in the current phase of [tasks.md](tasks.md). Mark it `[~]`.
2. **Size check before starting.** If it is `L`, split it in tasks.md first. If the time or AI usage budget left cannot finish it, do not start it. A half-finished slice is worse than none.
3. Work in **vertical slices**, in this order:
   1. **contract**: edit `contracts/openapi.yaml`, then run `make gen`
   2. **db**: add a migration and sqlc queries if needed, then run `make gen`
   3. **api**: handler → service → store in the owning module, with tests
   4. **web**: the screen or portal piece that uses it, through `@vernacular/api-client`
   5. **verify**: run it end to end (`make up`, `make dev`, then use it in the browser or with curl)

   Do not build several backend features and leave their frontends for later.
4. Run `make check`. It must pass.
5. For UI work, apply `/baseline-ui`. For any code task, run `/simplify` on the diff before committing.
6. Tick the box (and its sub-steps) in tasks.md, then make **one commit for the task** (see §6).

## 3. Reuse before writing

1. **Search the repo first.** Before writing a helper, component, query or type, search for an existing one (`rg <name>`). The duplication check (jscpd) fails CI above the threshold.
2. **Use a free, open-source library before writing our own.** Before adopting one, check its licence (MIT, Apache-2.0, BSD or ISC preferred; nothing non-commercial) and its maintenance (a release within the last 12 months, open issues answered). Record the choice in a short ADR if it is structural.
3. **Write from scratch only** what is specific to Zambian languages (stemmers, orthography rules, agreement scoring) or when nothing suitable exists. Say why in the PR.
4. **Extract on the second use, not the first.** Shared Go code goes to `apps/api/internal/platform`; shared TS goes to `packages/*`. A one-off does not need an abstraction.
5. **Generated code is never edited:** `internal/gen`, `*/dbgen`, `packages/api-client/src/schema.d.ts`. Change the source and run `make gen`.

## 4. Structure and conventions

### Go (`apps/api`)
- Layout: `cmd/<binary>` (wiring only), `internal/platform/*` (cross-cutting), `internal/<module>/` (domain).
- Inside a module: `handler.go` (implements the generated strict interface; maps types) → `service.go` (rules; no HTTP, no SQL) → `store.go` (sqlc and transactions).
- Modules call each other only through **service interfaces**, never through each other's stores or tables.
- Accept interfaces, return structs. Define each interface in the package that *uses* it.
- Errors: wrap with `fmt.Errorf("context: %w", err)`. Map domain errors to HTTP in one place (`platform/httpx`).
- Pass `context.Context` first on every call that does I/O. Do not keep global state.
- Tests: table-driven, `*_test.go` beside the code. Services get unit tests with fakes; stores get integration tests against Postgres.
- Text: run every input through `textnorm` (NFC). Never compare raw user strings.

### Web (`apps/web`)
- `app/(portal)` holds public pages and `app/(tools)` holds signed-in pages. Route files stay thin and compose components from `features/<domain>`.
- `components/ui` holds only shadcn primitives. Product components go in `features/<domain>/components`.
- Data comes from `@vernacular/api-client` through `lib/api.ts`. No hand-written `fetch` calls to the API, and no hand-written API types.
- Server Components by default; add `"use client"` only where interactivity needs it.
- Tailwind classes plus the shadcn tokens; no new colour values outside the theme.
- Accessibility: labelled controls, keyboard reachable, visible focus. The review screen must work on a phone.

### Database
- goose migrations in `apps/api/db/migrations`, named `NNNNN_description.sql`, each with `-- +goose Up` and `-- +goose Down`. Never edit a migration that has been merged; add a new one.
- Queries live in `apps/api/db/queries/<module>.sql`, with sqlc naming (`-- name: ListLanguages :many`).
- Every table has `created_at` and `updated_at`. Changes to reviewable items write to `audit_log` (NFR-8).

### Naming
- Language codes are ISO 639-3 (`bem`, `nya`, `toi`, `loz`, `kqn`, `lun`, `lue`, `eng`).
- IDs exposed in the API are stable strings and never change (PRV-2).
- JSON fields are `snake_case`; Go names are idiomatic Go; TS names are camelCase except for API types.

## 5. Definition of done

- [ ] The slice works end to end, and you ran it.
- [ ] Tests cover the new behaviour, including the failure cases.
- [ ] `make check` passes: lint, typecheck, tests, duplication check and generated-code drift.
- [ ] The contract, the generated code and the docs agree.
- [ ] The tasks.md box is ticked, and any new follow-up tasks are added.
- [ ] No secrets, personal data or unlicensed data is committed.

## 6. Commits

- **One task, one commit**, using Conventional Commits: `feat(api): …`, `feat(web): …`, `fix: …`, `docs: …`, `chore: …`, `refactor: …`, `test: …`.
- The subject is imperative and at most 72 characters. The body says *why* when that is not obvious, and references the task ID (`Task: P1-API-02`).
- **No AI co-author or attribution trailers** in commit messages.
- Do not commit `.env` files, keys or data dumps.

## 7. Commands

| Command | What it does |
|---|---|
| `make up` / `make down` | Start or stop Postgres (docker compose) |
| `make migrate` | Apply migrations |
| `make gen` | Regenerate Go server code, sqlc code and the TS client from their sources |
| `make dev` | Run the API and the web app |
| `make test` | Go and web tests |
| `make lint` | golangci-lint, ESLint and Prettier check |
| `make check` | Everything CI runs: lint, typecheck, test, duplication check, generated-code drift |
