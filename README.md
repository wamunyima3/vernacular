# Vernacular

**An open API of human-verified words, phrases and translations for Zambian languages.** Bemba and Nyanja come first, followed by Tonga, Lozi, Kaonde, Lunda and Luvale.

Developers call the API to add local language support to apps, government portals, SMS and USSD services. Every item the API returns says where it came from and who checked it. Nothing is marked `verified` until two native speakers have accepted it, independently. AI translations from Claude and Google are used to draft candidates for reviewers, but they stay labelled `machine` until people accept them.

> **Status:** Phase S (repository setup) is complete. Phase 0 (foundations: licences, schema, style guides) is next. See [tasks.md](tasks.md).

## Documentation

| Document                                     | What it covers                                                                         |
| -------------------------------------------- | -------------------------------------------------------------------------------------- |
| [docs/requirements.md](docs/requirements.md) | Product requirements, with IDs such as `LEX-4` and `DQ-1` (snapshot of the design doc) |
| [docs/user-stories.md](docs/user-stories.md) | User stories per persona, with acceptance criteria                                     |
| [docs/architecture.md](docs/architecture.md) | System shape, modules, data model, trust model, contract-first flow                    |
| [docs/adr/](docs/adr/)                       | Architecture decision records                                                          |
| [tasks.md](tasks.md)                         | The task tracker: every piece of work, by phase                                        |
| [AGENTS.md](AGENTS.md)                       | Rules for developers and AI agents: conventions, workflow, definition of done          |
| [CONTRIBUTING.md](CONTRIBUTING.md)           | How to contribute code or language data                                                |

## Repository layout

```
contracts/openapi.yaml   API contract: the source of truth for server and client code
apps/api/                Go API (chi, pgx, sqlc, goose, oapi-codegen)
apps/web/                Next.js 16 developer portal and contributor/reviewer tools
packages/api-client/     Typed TS client generated from the contract
packages/eslint-config/  Shared lint rules
packages/tsconfig/       Shared TS config
deploy/                  docker compose for local services
docs/                    Requirements, stories, architecture, ADRs
```

## Quickstart

**Prerequisites:** Go 1.25 or newer (Go downloads the 1.26 toolchain the module asks for automatically), Node 22+, pnpm 10, Docker, and [golangci-lint](https://golangci-lint.run) v2 (`brew install golangci-lint`).

```bash
cp .env.example .env
make setup      # pnpm install, Go modules, git hooks
make up         # start Postgres
make migrate    # create tables and seed languages
make dev        # API on :8080, web on :3000
```

Then open http://localhost:3000, or call the API directly:

```bash
curl localhost:8080/v1/languages
```

Codegen tools (oapi-codegen, sqlc, goose) are pinned in `apps/api/tools/go.mod` and run through `go tool`, so there is nothing else to install.

## Everyday commands

| Command      | What it does                                                                               |
| ------------ | ------------------------------------------------------------------------------------------ |
| `make help`  | List all targets                                                                           |
| `make gen`   | Regenerate Go server code, sqlc code and the TS client after changing the contract or SQL  |
| `make test`  | Go and web tests                                                                           |
| `make lint`  | golangci-lint, ESLint and the Prettier check                                               |
| `make check` | Everything CI runs: generated-code drift, lint, typecheck, tests, duplication check, build |

## How work happens

1. Pick the next unchecked task in [tasks.md](tasks.md), and check its size first.
2. Build it as a vertical slice: contract → db → api → web → verify.
3. Run `make check`, tick the task, and make one Conventional Commit per task.

The full rules are in [AGENTS.md](AGENTS.md).

## Licence

Code: [MIT](LICENSE). The dataset will be published under its own open licence (CC BY 4.0 is proposed). Only sources that allow open redistribution go into it.
