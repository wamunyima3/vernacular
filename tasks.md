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
- [x] **P0-CNT-03** Core concept list (~4,000) → LEX-7: 4,073 core concepts after `make seed`; senior reviewers should prune it in Phase 1
  - [x] **P0-CNT-03a** Seed concepts from the English side of GATITOS en→bem, with a repeatable import command (S): `make seed` gives 3,977 core concepts and 2,848 Bemba entries, all `sourced`
  - [x] **P0-CNT-03b** Add Zambian everyday terms by domain (health, farming, money, government, family) (M): `data/concepts/zambian-everyday.csv`, 255 concepts

### Engineering

- [x] **P0-ENG-01** Full core schema as goose migrations: language, concept, entry, variant, sense, example, phrase, source, review, release, audit log (M) → LEX-3, LEX-6, LEX-7, PRV-2, NFR-8
- [x] **P0-ENG-02** Shared trust fields (`status`, `sources`, `reviews`, `release`) as one OpenAPI component, reused by every item schema (S) → PRV-1, ADR-0004
- [x] **P0-ENG-03** Uniform error body (`code`, `message`, `docs_url`) and error middleware (S) → §10 errors
- [x] **P0-ENG-04** Text normalisation package: NFC, case and accent folding (`golang.org/x/text`) with table tests (S) → NFR-7
- [x] **P0-ENG-05** Evaluate Better Auth vs Auth.js against the roles developer, contributor, reviewer, senior reviewer and admin; write ADR (S)
- [x] **P0-ENG-06** Auth slice: sign up, sign in and roles (M) → US-01, US-18
  - [x] api: users and roles tables, session validation middleware
  - [x] web: sign-up and sign-in screens
  - [x] test: role guard tests
- [x] **P0-ENG-07** UI designs for portal and review tools; run `/baseline-ui` on them (M)
- [x] **P0-ENG-08** ~~Throwaway lookup prototype against a few hundred hand-checked entries~~ Replaced: no hand-checked entries exist until reviewers are recruited (P0-OUT-07), so the real search slice (P1-API-02) is built on GATITOS `sourced` data instead (M)
- [x] **P0-ENG-09** Evaluate Fumadocs vs Nextra, and Scalar for the API reference; write ADR (S) → DEV-2

---

## Phase 1: Bemba and Nyanja lexicon (~3 months)

Done when: 4,000 verified concepts exist in both languages, and the API is open to a closed group of developers.

### Sources and ingestion

- [x] **P1-SRC-01** Source register slice (M) → DQ-1, US-25
  - [x] api: CRUD for sources, permission document upload (S3-compatible)
  - [x] web: `(tools)/sources` screen
  - [x] test: an ingest is refused for a source without a register row
- [x] **P1-SRC-02** Ingestion framework: `cmd/ingest` with one parser per source behind a shared interface (M) → US-30
- [x] **P1-SRC-03** Normaliser: orthography rules, duplicate merge, concept linking (M) → LEX-3, LEX-7: style-guide spellings with source spellings kept as variants; concepts link by normalised English gloss
- [x] **P1-SRC-04** Rule-based stemmer for Bemba (M) → LEX-4
- [x] **P1-SRC-05** Rule-based stemmer for Nyanja (M) → LEX-4: rules from the style guide; tune once Nyanja data arrives

### AI-assisted candidates

- [x] **P1-AI-01** River job: generate Claude candidates per concept and phrase, stored as `machine` with provider, model version and date (M) → AI-1, AI-2, US-27
  - Concepts done (`make candidates LANG=bem`, `make dev-worker`); phrases follow once the phrase list exists (P1-API phrases). A live run needs `ANTHROPIC_API_KEY` (ADR-0009).
- [x] **P1-AI-02** River job: generate Google Translate candidates (bem, ny only) (S) → AI-1, AI-2
- [x] **P1-AI-03** Agreement scoring and review priority (M) → AI-3, US-21

- [x] **P1-AI-04** Phrase candidates: P1-AI-01 generated words only, so the phrase packs are still empty (found during P2-05). Split:
  - [x] **P1-AI-04a** River job: Claude and Google propose each phrase-list line in a language, stored as `machine` phrases; identical proposals agree on one item (M) → AI-1, AI-2, PHR-1
  - [x] **P1-AI-04b** The review queue serves phrase candidates: cards, decisions, promotion and "write my own" for phrases (M) → CON-2, US-20

### Review

- [x] **P1-AI-03a** Apply the /simplify findings on P1-AI-02/03 (S): agreement and priority are now the `review_agreement` view, read by the queue; restem reads active languages; provider source names live in `candidates`

- [x] **P1-REV-00** Give `items` and `item_sources` one owner (`internal/items`) with a store that lexicon calls, and have `sources` depend on a lexicon writer interface instead of `lexicon.NewStore` (S) → architecture §5
- [x] **P1-REV-01** Contribute slice: propose entry, sense and example (M) → CON-1, US-18
  - [x] api: `POST /v1/contributions`, stored as `sourced` citing "Community contributors" and the author
  - [x] web: `(tools)/contribute`, with the trust badge on the receipt
- [x] **P1-REV-01b** Corrections: propose a change to an existing item, linked to it and kept in its history (M) → CON-1, US-18, PRV-2. Needs item reads (P1-API) and promotion rules (P1-REV-03)
  - [x] **P1-REV-01b1** api + web: a contribution can name the meaning it `corrects`; it is a new candidate linked by `items.corrects`, and the original stays served. "Suggest a correction" on each meaning opens the contribute form prefilled (S)
  - [x] **P1-REV-01b2** Reviewers see a correction beside the meaning it corrects; when the correction is verified the original is retired (kept, with its history, never served) (M)
- [x] **P1-REV-02** Review queue slice, blind and mobile-first → CON-2, CON-3, DQ-4, AI-4, US-20, US-21
  - [x] **P1-REV-02a** api: `GET /v1/review/next` gives one concept with every pending candidate, no provider names, quick confirms first; `POST /v1/reviews` records accepts and rejects. Only approved languages, never the reviewer's own or already-reviewed items. Decisions only; promotion is P1-REV-03 (M)
  - [x] **P1-REV-02b** web: `(tools)/review`, one card per screen, fixed action bar, reject reasons, skip, empty state (M)
  - [x] **P1-REV-02c** Undo a decision for 5 seconds (S)
- [x] **P1-REV-03** Promotion and corrections → DQ-2, PRV-2, US-22
  - [x] **P1-REV-03a** One promotion function: two independent accepts make an item `verified`, one makes it `reviewed`, disagreement makes it `disputed`, two rejects make it `rejected` (new status, never returned). Decide and undo both recompute through it, so a promoting accept can be undone (M)
  - [x] **P1-REV-03b** "Write my own": a reviewer's correction becomes a new candidate for others to review; the candidates it replaces keep their history (M)
- [x] **P1-REV-04c** Disputed phrases reach senior reviewers: `/disputes` lists senses only, so a disputed phrase waits there unseen (S) → CON-4, DQ-8
- [x] **P1-REV-04** Disputes slice for senior reviewers; reports and senior resolutions become inputs to `review.Status`, so an undo cannot overwrite them (M). Reports join with P1-API-10 → CON-4, DQ-8, US-23
- [x] **P1-REV-05** Reviewer qualification test set and agreement-rate tracking (M) → DQ-3, DQ-7, US-24
  - [x] **P1-REV-05a** Agreement rates: each reviewer's share of decisions that match settled outcomes, shown to senior reviewers; below the threshold a reviewer loses approval until a senior marks them retrained (M)
  - [x] **P1-REV-05b** Qualification: a short test set from verified items for new reviewers, confirmed by a senior reviewer (M). Needs verified items in the language
    - [x] **P1-REV-05b1** api: `POST /v1/qualifications` draws ten pick-the-word questions from verified words (answers kept server side), `POST …/answers` scores once (8 of 10 passes), seniors list passes and confirm one, which approves the taker as a reviewer. Once a day until passed (M)
    - [x] **P1-REV-05b2** web: take the test in `(tools)/qualify`, see the score; seniors confirm passes on the reviewers page (M)
- [x] **P1-REV-06** Admin: approve reviewers, acceptance rates per provider and language (M) → AI-5, US-26
- [x] **P1-REV-07** Public contributor credit page (S) → CON-5, US-19
- [x] **P1-REV-08** Rename the review API's `sense_id` to `item_id` (decisions, candidates, reviews, corrections) now that phrases are reviewed too, and give the review package one candidate type for senses and phrases (S)

### Public API (each task is a slice: contract → api → portal "Try it" → test)

- [x] **P1-API-01** `GET /v1/languages` with live entry counts (S) → US-03
- [x] **P1-API-02** `GET /v1/entries/search` with variant and root matching (M) → LEX-1, LEX-4, US-03, US-05
- [x] **P1-API-03** `GET /v1/entries/{id}` (S) → LEX-1, LEX-3, LEX-5, US-16
- [x] **P1-API-04** `GET /v1/concepts/{id}` (S) → LEX-7, US-04
- [x] **P1-API-05** `GET /v1/lookup` in both directions with English (M) → LEX-2, US-04
- [x] **P1-API-06** `min_status` and `variety` filters, verified by default (S) → TRN-4, AI-6, LEX-6, US-06, US-07
- [x] **P1-API-07** `GET /v1/phrases` packs (M) → PHR-1, US-08
- [x] **P1-API-07b** English source list for the six packs (greetings, numbers, dates, interface labels, health, money), loaded as phrase concepts, so P1-AI-01 can generate phrase candidates and reviewers can check them (S) → PHR-1, AI-1
- [x] **P1-API-08** `POST /v1/detect` (evaluate existing language ID models first) (M) → TRN-1, US-09
- [x] **P1-API-09** Exact-match translation against verified phrases (S) → TRN-2, US-10
- [x] **P1-API-10** `POST /v1/reports` and the report inbox for reviewers; a confirmed report is an input to `review.Status` (→ `disputed`) and shows on `/disputes` (M) → PRV-3, US-15
  - The portal's "Report a problem" button comes with the search results in P1-WEB-02, which calls `POST /v1/reports`.
- [x] **P1-API-11** API keys slice: create and revoke keys in the portal; keys hashed at rest (M) → DEV-1, NFR-5, US-01
- [x] **P1-API-12** Rate limiting per key and per IP, with rate-limit headers; derive client IP only from trusted proxy headers, not chi RealIP (S) → DEV-3, NFR-5, US-13
- [x] **P1-API-13** Compression, `ETag` tied to the release, cache headers (S) → NFR-3, US-11
  - No releases exist yet, so the `ETag` is a hash of the body; P2 release work ties it to the release number.

### Portal

- [x] **P1-WEB-01** Home with live verified counts (S)
- [x] **P1-WEB-02a** Try it: `/try` search with the formatted result and its JSON side by side (tabs on phones), and a search box on the home page (S) → US-02
  - Signed-out visitors search with the portal's key (`WEB_API_KEY`); without one they are asked to sign in.
- [x] **P1-WEB-02b** Entry page `/entries/[id]` with senses, examples, variants and sources, linked from the search results, and a "Report a problem" dialog on each meaning; the portal header shows "Account" when signed in (M) → US-16, US-15
- [x] **P1-WEB-03** Interactive API reference from the contract (S) → DEV-2
- [x] **P1-WEB-04** Docs: concepts, verification levels, varieties, licences (M) → DEV-2
- [x] **P1-WEB-05** Tutorial and copy-ready examples (JS, Go, Python, curl) (M) → DEV-2

### Operations

- [x] **P1-OPS-01** Container images and a cheap single-server deploy recipe (M) → NFR-10
- [x] **P1-OPS-02** p95 latency measurement on read endpoints (S) → NFR-1. Measured locally on the seed data: p95 about 9 ms on search, lookup and phrases; re-measure on the server once real data is in

---

## Phase 2: Tonga and Lozi, plus polish (~3 months)

Done when: four languages are live and the dictionary and phrase API has launched publicly.

- [ ] **P2-01** Add Tonga and Lozi: style guides, stemmers, AI candidates (Claude only) (L, split)
  - [x] **P2-01a** Tonga (`toi`): draft style guide, spelling rules (`ch` → `c`, `ŋ` → `ng'`), stemmer, Claude prompt (S). Google Translate does not offer it
  - [x] **P2-01b** Lozi (`loz`): draft style guide, spelling rules (`ch` → `c`, `ng'`/`ŋ` → `ñ`), stemmer, Claude prompt (S)
  - [ ] **P2-01c** Activate Tonga and Lozi once each style guide is signed off by two senior reviewers 👤
- [x] **P2-02** Release builder job and `GET /v1/releases` (M) → PRV-4, US-11, US-29. Split:
  - [x] **P2-02a** Build a release (admin, `POST /v1/releases`) that freezes the verified items with counts and changes against the previous release; `GET /v1/releases`; the portal releases page (M). Built in the request, in one transaction; move it to a River job if a build outgrows a request
  - [x] **P2-02b** Pinned reads (ADR-0012). Split:
    - [x] **P2-02b1** The builder materialises each release as a schema (`materialise_release`) (S)
    - [x] **P2-02b2** `release=` on the read endpoints, through a per-request connection on the release's `search_path`; `trust.release` filled; the `ETag` tied to the release (M) → US-11, NFR-3
- [x] **P2-03** Open dataset export in open formats and the download page (M) → NFR-9, US-17
- [x] **P2-04** `GET /v1/phrases/export` in JSON, PO and XLIFF (M) → PHR-3, US-12
- [x] **P2-05** Phrase placeholders with grammar notes (M) → PHR-2, US-32
- [x] **P2-06** JavaScript and Python client libraries (M) → DEV-4, US-14
  - [x] **P2-06a** `@vernacular/client` (`packages/client`): openapi-fetch typed from the contract, the key sent for you, named methods for the public reads that throw `VernacularError`, built to `dist` (S)
  - [x] **P2-06b** `vernacular-client` (`clients/python`): openapi-python-client output in `vernacular/_api` from `make gen`, and a `Vernacular` class with the same named methods raising `VernacularError`; pytest and ruff in CI through uv (M)
  - [ ] **P2-06c** Publish both to npm and PyPI under the project's accounts 👤
- [x] **P2-07** Usage dashboard per key (M) → DEV-5, US-13
  - [x] **P2-07a** api: calls per key per UTC day, counted in memory and flushed every 30 s into `key_usage` (calls, errors, 429s); `GET /v1/keys/{id}/usage` for the key's owner, with the per-minute limit (M)
  - [x] **P2-07b** web: today and 30-day calls per key on the account page (totals now in `GET /v1/keys`), and `/account/keys/[id]` with daily calls, errors, 429s and the per-minute limit (S)
- [ ] **P2-08** Pronunciation audio with consent records (L, split) → LEX-8, NFR-6, US-33
  - [x] **P2-08a** api: speakers and their consent records: an admin records a speaker's signed consent form (stored like permission documents), a speaker sees and withdraws their consent; everything behind `AUDIO_ENABLED`, off by default (M)
  - [x] **P2-08b** api: clips: a speaker whose consent form names the entry's language records it (Ogg, WebM, MP4, MP3 or WAV, at most 1 MB; again replaces); `GET /v1/audio?entry_id=`, `GET`/`DELETE /v1/audio/{id}` (speaker or admin); withdrawing hides the clips at once, then deletes them and their files under `speakers/{id}/` (M)
  - [x] **P2-08c** web: entry pages play an entry's clips (streamed through `/audio/[id]`) and give a consenting speaker of the language a MediaRecorder recorder; a speaker deletes their own clip, and the account page shows their consent with a withdraw button (M)
  - [ ] **P2-08d** Take legal advice, register as a data controller, decide where audio is hosted, then set `AUDIO_ENABLED` 👤
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
