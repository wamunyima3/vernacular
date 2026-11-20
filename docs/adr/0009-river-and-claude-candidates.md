# ADR-0009: River for background jobs; Claude candidates per concept

- **Status:** Accepted
- **Date:** 2026-11-20
- **Task:** P1-AI-01

## Context

AI candidates (AI-1) mean one model call per concept and language: about 4,000 concepts × 2 languages to start. Calls can take minutes, fail transiently, or be declined, and must not be lost or repeated twice when the run is restarted. The architecture already names a worker running River (§2).

| Criterion           | River 0.51                              | A cron loop in our own code |
| ------------------- | --------------------------------------- | --------------------------- |
| Licence             | MPL-2.0 (file-level; we do not fork it) | ours                        |
| Last release        | 2026-11, monthly                        | n/a                         |
| Queue storage       | The Postgres we already run             | a table we would design     |
| Retries, uniqueness | Built in, per job                       | to be written               |
| Enqueue in a tx     | Yes (`InsertManyTx`)                    | to be written               |

MPL-2.0 is not in the preferred list (AGENTS.md §3), but it is not non-commercial and it only binds changes to River's own files, so it is compatible with an MIT codebase that depends on it unmodified.

## Decision

- **River** runs jobs in `cmd/worker`. Its tables come from its own migrations (`river migrate-up`, run by `make migrate`), not copied into goose.
- **One job per concept and language**, unique by arguments, so `make candidates LANG=bem` can be re-run and only queues what is missing. A refusal cancels the job; other errors retry (up to 5 attempts).
- **Claude Opus 5.5 through the Messages API**, not the Batches API: server-side refusal fallbacks are not available on batches, and jobs give us per-concept retries. The system prompt carries the language's style guide and is cached, so each call adds only the concept. Structured outputs constrain the answer to the `part_of_speech` enum.
- Candidates are written through `lexicon.Add` with a `machine` origin citing the `Claude` source row and the model that answered. When a word already exists from a licensed source, its status is kept and Claude is added to its sources, which is the agreement signal for P1-AI-03.

## Consequences

- Running candidates needs `ANTHROPIC_API_KEY` and spends money: about 8,000 calls for both languages.
- If cost matters more than per-concept retries, the Batches API halves the price; revisit then.
- The worker holds no transaction open during a model call; reads and writes are separate short transactions.
