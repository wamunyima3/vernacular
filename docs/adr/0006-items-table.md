# ADR-0006: One `items` table for everything reviewable

- **Status:** Accepted
- **Date:** 2026-11-20
- **Task:** P0-ENG-01

## Context

Entries, senses, examples and phrases are all reviewed, sourced, released and audited the same way (PRV-1, PRV-2, DQ-2). The architecture left open whether reviews point at items through a polymorphic `(item_type, item_id)` pair, or whether each type gets its own review table.

## Decision

Every reviewable row has a row in `items` with the same UUID: `items(id, kind, language_code, status, author_id)`. `entries.id`, `senses.id`, `examples.id` and `phrases.id` are foreign keys to `items.id`. `reviews`, `item_sources` and `release_items` reference `items`.

Two triggers apply to every table through `track(table)`: `updated_at` maintenance, and a row in `audit_log` with before/after JSON and the actor from `app.actor_id` (NFR-8).

Search folding is one SQL function, `fold()`, used by generated columns and indexes. The Go `textnorm.Fold` (P0-ENG-04) must match it, and a test checks that.

## Consequences

- Reviews, sources and releases have real foreign keys, and one query serves the whole review queue.
- Trust status lives in one column, so promotion (P1-REV-03) has one place to write.
- Creating an item takes two inserts. Stores do both in one transaction through a shared helper.
- The database does not stop an `entries` row pointing at an `items` row of another kind. The shared insert helper is the only writer.
- User IDs are text without foreign keys until the auth task adds the users table.
