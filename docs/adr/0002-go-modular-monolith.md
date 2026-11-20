# 0002. Go modular monolith

- Status: Accepted
- Date: 2026-11-20

## Context

The requirements (§9, NFR-10) call for a single modest server, a small team and mostly read traffic. Microservices would add deployment and consistency cost without benefit at this size.

## Decision

Build one Go codebase with several thin binaries (`api`, `worker`, `ingest`) that share `internal/` packages. Domain modules live in `internal/<module>` with a fixed handler → service → store layering. Modules talk to each other only through service interfaces. The Phase 3 translation service is the one planned exception, because it has a different runtime profile.

## Consequences

- One deploy, one database, transactions across modules are possible.
- Module boundaries must be enforced by review and lint rules, not by the network.
- Splitting out a module later means swapping an interface implementation for a client.
