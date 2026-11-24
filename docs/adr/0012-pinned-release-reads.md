# 0012. Pinned reads through a schema per release

- Status: Accepted
- Date: 2026-11-24

## Context

US-11 lets a caller pin a release (`release=2027.1`) so the data does not change under their app. Every read endpoint (search, entries, lookup, concepts, phrases, translate) would then have to serve the items as they were frozen. The release builder (P2-02a) already keeps a JSON snapshot of each item, but serving reads from JSON would need a second copy of every read query, in SQL that differs from the live one, across the lexicon, phrases and items modules.

## Decision

The release builder also materialises the release as a PostgreSQL schema, `release_<year>_<n>`, holding copies of the released rows of `items`, `entries`, `senses`, `examples`, `phrases`, `variants` and `item_sources`, with the live tables' indexes. A SQL function, `materialise_release`, does this inside the build transaction.

A pinned request takes one pooled connection for its whole lifetime and sets `search_path` to the release schema, then `public`. The existing, unqualified queries then read the frozen rows; tables that are not copied (concepts, sources, languages) fall through to `public`. Stores receive a database handle that uses the request's connection when there is one.

Without `release`, reads use the live verified data, as before. `release=latest` names the newest release. This departs from the requirements' default of `latest`: until the first release is built, a `latest` default would return nothing.

## Consequences

- Pinned reads reuse every read query unchanged, and pinned search keeps its indexes.
- Each release costs a copy of the verified rows. Old releases can be dropped by schema once their support window ends.
- Concept glosses, source details and review counts are live, not frozen. They describe items rather than being the data a caller pins.
- A pinned request holds one connection while it runs, as a transaction would.
