# Vernacular: task tracker

The only list of work for this repo. Tick a box in the same commit that finishes the task.

- **Requirements:** [docs/requirements.md](docs/requirements.md) (IDs such as `LEX-4`, `DQ-1`)
- **User stories:** [docs/user-stories.md](docs/user-stories.md) (IDs such as `US-05`)
- **Architecture:** [docs/architecture.md](docs/architecture.md)
- **Working rules:** [AGENTS.md](AGENTS.md)

## How to use this file

1. Take the **first unchecked task** in the current phase whose dependencies are done. Phases run in order; tasks inside a phase mostly run in order too.
2. **Check the size before starting.** `S` takes under an hour, `M` takes one focused session, and `L` is too big. Split an `L` into `S` or `M` tasks here first, then start. If the usage budget left in the session cannot cover the task, stop rather than leave it half done.
3. A code task is a **vertical slice**. Work through its sub-steps in order: contract → api → web → test. It is done only when the whole slice works end to end.
4. Finish against the definition of done in [AGENTS.md](AGENTS.md#definition-of-done), tick the box, and commit (one task, one commit).
5. If you find work that is not listed, add it as a task in the right phase before doing it.

**Legend:** `[ ]` todo · `[x]` done · `[~]` in progress (one at a time) · `(S/M/L)` size · `→` refs to requirements and stories · 👤 owner task (not code; Wamunyima)

---

## Phase S: Repository setup

Goal: a documented, runnable skeleton with quality gates, so feature work can start.

- [x] **S-01** Requirements snapshot and task tracker (S) → all
- [x] **S-02** User stories mapped to requirement IDs (M) → all Must items
- [x] **S-03** Architecture document and ADRs 0001–0005 (M) → §9, §10
- [x] **S-04** AGENTS.md, CLAUDE.md and CONTRIBUTING.md (S)
- [x] **S-05** Monorepo tooling: pnpm workspace, Turborepo, Makefile, docker compose, MIT licence (S) → NFR-10
- [x] **S-06** Go API skeleton (M) → NFR-4, NFR-7
  - [x] contract: `contracts/openapi.yaml` with `/healthz` and `GET /v1/languages`
  - [x] api: chi server, config, slog, pgx pool, oapi-codegen strict server
  - [x] db: goose migration for `languages` (seed bem, nya); sqlc query
  - [x] test: handler test for `/v1/languages`
- [x] **S-07** Next.js web skeleton (M) → DEV-2
  - [x] web: Next.js 16 + Tailwind + shadcn/ui, `(portal)` and `(tools)` route groups
  - [x] packages: `@vernacular/api-client` generated from the contract
  - [x] slice: portal home lists the supported languages from the live API
- [x] **S-08** Quality gates and CI: golangci-lint, ESLint, Prettier, jscpd, lefthook, generated-code drift check, GitHub Actions (M)
- [x] **S-09** `/simplify` pass over the scaffold (S)
- [x] **S-10** Rewrite README (S)

---

## Phase 0: Foundations (4–6 weeks)

Done when: at least one licensed source per language is signed, and the schema and style guides are agreed.

### Outreach and legal 👤

- [ ] **P0-OUT-01** Write to the UNZA speech lab about partnership and licence terms for Zambezi Voice and BIG-C 👤
- [ ] **P0-OUT-02** Write to the Bible Society of Zambia about text licences for all seven languages, including the right to republish extracted words and sentences 👤
- [ ] **P0-OUT-03** Contact the Kitwe Online and Chichewa Dictionary authors about licensing their word lists 👤
- [ ] **P0-OUT-04** Get the full BongoHive ecosystem report and add what it lists to the requirements doc, §3 👤
- [ ] **P0-OUT-05** Get legal advice on the copyright status of the pre-1960 dictionaries and on Data Protection Act duties 👤
- [ ] **P0-OUT-06** Obtain the Ministry of Education orthography for each language 👤
- [ ] **P0-OUT-07** Find two senior reviewers each for Bemba and Nyanja 👤
- [ ] **P0-OUT-08** Choose the data licence (CC BY 4.0 proposed) and file the Claude and Google output permissions in the source register 👤 → §12
- [ ] **P0-OUT-09** Recruitment plan for volunteer reviewers (Byte Builders, UNZA students, teachers, churches, radio) 👤

### Content

- [x] **P0-CNT-01** Download Google SMOL and GATITOS and list exactly which Zambian languages they cover (S) → [findings](docs/sources/smol-gatitos.md)
- [x] **P0-CNT-02** Draft the Bemba and Nyanja style guides in `docs/style-guides/` (M) → DQ-5
- [ ] **P0-CNT-03** Core concept list (~4,000) → LEX-7
  - [ ] **P0-CNT-03a** Seed concepts from the English side of GATITOS en→bem, with a repeatable import command (S)
  - [ ] **P0-CNT-03b** Add Zambian everyday terms by domain (health, farming, money, government, family) (M)

### Engineering

- [x] **P0-ENG-01** Full core schema as goose migrations: language, concept, entry, variant, sense, example, phrase, source, review, release, audit log (M) → LEX-3, LEX-6, LEX-7, PRV-2, NFR-8
- [x] **P0-ENG-02** Shared trust fields (`status`, `sources`, `reviews`, `release`) as one OpenAPI component, reused by every item schema (S) → PRV-1, ADR-0004
- [x] **P0-ENG-03** Uniform error body (`code`, `message`, `docs_url`) and error middleware (S) → §10 errors
- [ ] **P0-ENG-04** Text normalisation package: NFC, case and accent folding (`golang.org/x/text`) with table tests (S) → NFR-7
- [ ] **P0-ENG-05** Evaluate Better Auth vs Auth.js against the roles developer, contributor, reviewer, senior reviewer and admin; write ADR (S)
- [ ] **P0-ENG-06** Auth slice: sign up, sign in and roles (M) → US-01, US-18
  - [ ] api: users and roles tables, session validation middleware
  - [ ] web: sign-up and sign-in screens
  - [ ] test: role guard tests
- [ ] **P0-ENG-07** UI designs for portal and review tools; run `/baseline-ui` on them (M)
- [ ] **P0-ENG-08** Throwaway lookup prototype against a few hundred hand-checked entries (M) → first ten tasks, #10
- [ ] **P0-ENG-09** Evaluate Fumadocs vs Nextra, and Scalar for the API reference; write ADR (S) → DEV-2

---

## Phase 1: Bemba and Nyanja lexicon (~3 months)

Done when: 4,000 verified concepts exist in both languages, and the API is open to a closed group of developers.

### Sources and ingestion

- [ ] **P1-SRC-01** Source register slice (M) → DQ-1, US-25
  - [ ] api: CRUD for sources, permission document upload (S3-compatible)
  - [ ] web: `(tools)/sources` screen
  - [ ] test: an ingest is refused for a source without a register row
- [ ] **P1-SRC-02** Ingestion framework: `cmd/ingest` with one parser per source behind a shared interface (M) → US-30
- [ ] **P1-SRC-03** Normaliser: orthography rules, duplicate merge, concept linking (M) → LEX-3, LEX-7
- [ ] **P1-SRC-04** Rule-based stemmer for Bemba (M) → LEX-4
- [ ] **P1-SRC-05** Rule-based stemmer for Nyanja (M) → LEX-4

### AI-assisted candidates

- [ ] **P1-AI-01** River job: generate Claude candidates per concept and phrase, stored as `machine` with provider, model version and date (M) → AI-1, AI-2, US-27
- [ ] **P1-AI-02** River job: generate Google Translate candidates (bem, ny only) (S) → AI-1, AI-2
- [ ] **P1-AI-03** Agreement scoring and review priority (M) → AI-3, US-21

### Review

- [ ] **P1-REV-01** Contribute slice: propose entry, sense, example or correction (M) → CON-1, US-18
- [ ] **P1-REV-02** Review queue slice, blind and mobile-first (L, split before starting) → CON-2, CON-3, DQ-4, AI-4, US-20
- [ ] **P1-REV-03** Promotion rules: two independent accepts make an item `verified`; a correction keeps its history (M) → DQ-2, PRV-2, US-22
- [ ] **P1-REV-04** Disputes slice for senior reviewers (M) → CON-4, DQ-8, US-23
- [ ] **P1-REV-05** Reviewer qualification test set and agreement-rate tracking (M) → DQ-3, DQ-7, US-24
- [ ] **P1-REV-06** Admin: approve reviewers, acceptance rates per provider and language (M) → AI-5, US-26
- [ ] **P1-REV-07** Public contributor credit page (S) → CON-5, US-19

### Public API (each task is a slice: contract → api → portal "Try it" → test)

- [ ] **P1-API-01** `GET /v1/languages` with live entry counts (S) → US-03
- [ ] **P1-API-02** `GET /v1/entries/search` with variant and root matching (M) → LEX-1, LEX-4, US-03, US-05
- [ ] **P1-API-03** `GET /v1/entries/{id}` (S) → LEX-1, LEX-3, LEX-5, US-16
- [ ] **P1-API-04** `GET /v1/concepts/{id}` (S) → LEX-7, US-04
- [ ] **P1-API-05** `GET /v1/lookup` in both directions with English (M) → LEX-2, US-04
- [ ] **P1-API-06** `min_status` and `variety` filters, verified by default (S) → TRN-4, AI-6, LEX-6, US-06, US-07
- [ ] **P1-API-07** `GET /v1/phrases` packs (M) → PHR-1, US-08
- [ ] **P1-API-08** `POST /v1/detect` (evaluate existing language ID models first) (M) → TRN-1, US-09
- [ ] **P1-API-09** Exact-match translation against verified phrases (S) → TRN-2, US-10
- [ ] **P1-API-10** `POST /v1/reports` and the report inbox for reviewers (M) → PRV-3, US-15
- [ ] **P1-API-11** API keys slice: create and revoke keys in the portal; keys hashed at rest (M) → DEV-1, NFR-5, US-01
- [ ] **P1-API-12** Rate limiting per key and per IP, with rate-limit headers; derive client IP only from trusted proxy headers, not chi RealIP (S) → DEV-3, NFR-5, US-13
- [ ] **P1-API-13** Compression, `ETag` tied to the release, cache headers (S) → NFR-3, US-11

### Portal

- [ ] **P1-WEB-01** Home with live verified counts (S)
- [ ] **P1-WEB-02** Try it: search box with JSON beside the result (M) → US-02
- [ ] **P1-WEB-03** Interactive API reference from the contract (S) → DEV-2
- [ ] **P1-WEB-04** Docs: concepts, verification levels, varieties, licences (M) → DEV-2
- [ ] **P1-WEB-05** Tutorial and copy-ready examples (JS, Go, Python, curl) (M) → DEV-2

### Operations

- [ ] **P1-OPS-01** Container images and a cheap single-server deploy recipe (M) → NFR-10
- [ ] **P1-OPS-02** p95 latency measurement on read endpoints (S) → NFR-1

---

## Phase 2: Tonga and Lozi, plus polish (~3 months)

Done when: four languages are live and the dictionary and phrase API has launched publicly.

- [ ] **P2-01** Add Tonga and Lozi: style guides, stemmers, AI candidates (Claude only) (L, split)
- [ ] **P2-02** Release builder job and `GET /v1/releases` (M) → PRV-4, US-11, US-29
- [ ] **P2-03** Open dataset export in open formats and the download page (M) → NFR-9, US-17
- [ ] **P2-04** `GET /v1/phrases/export` in JSON, PO and XLIFF (M) → PHR-3, US-12
- [ ] **P2-05** Phrase placeholders with grammar notes (M) → PHR-2, US-32
- [ ] **P2-06** JavaScript and Python client libraries (M) → DEV-4, US-14
- [ ] **P2-07** Usage dashboard per key (M) → DEV-5, US-13
- [ ] **P2-08** Pronunciation audio with consent records (L, split) → LEX-8, NFR-6, US-33
- [ ] **P2-09** Example sentences for every sense in the core list (M) → LEX-5
- [ ] **P2-10** Quarterly 5% re-check sampling job (S) → DQ-6, US-28

---

## Phase 3: Machine translation beta (~3 months)

- [ ] **P3-01** Choose a base model whose licence allows commercial serving; write ADR 👤
- [ ] **P3-02** Python training pipeline in `ml/` (L, split) → AI-7
- [ ] **P3-03** Translation service (separate deployable) (M)
- [ ] **P3-04** `POST /v1/translate`: verified phrase first, then labelled machine output (M) → TRN-3, US-31
- [ ] **P3-05** Benchmark against Google and NLLB-200 on FLORES and publish the scores (M) → DQ-9, DQ-10

---

## Phase 4: Kaonde, Lunda, Luvale and speech (ongoing)

- [ ] **P4-01** Volunteer collection campaigns for Kaonde, Lunda and Luvale 👤
- [ ] **P4-02** Add each language through the same pipeline (L, split per language)
- [ ] **P4-03** Speech recognition and text to speech, if licences and funding allow (L, split)
