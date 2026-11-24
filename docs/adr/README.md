# Architecture decision records

One file per decision that is hard to reverse. Copy `0000-template.md`, take the next number, and link it from [architecture.md](../architecture.md) when it changes the structure.

| ADR                                         | Title                                           | Status   |
| ------------------------------------------- | ----------------------------------------------- | -------- |
| [0001](0001-monorepo.md)                    | Monorepo with pnpm workspaces and Turborepo     | Accepted |
| [0002](0002-go-modular-monolith.md)         | Go modular monolith                             | Accepted |
| [0003](0003-openapi-first.md)               | OpenAPI-first contract with generated code      | Accepted |
| [0004](0004-verification-levels.md)         | Verification levels and shared trust fields     | Accepted |
| [0005](0005-concept-hub.md)                 | Concepts as the hub between languages           | Accepted |
| [0006](0006-items-table.md)                 | One `items` table for everything reviewable     | Accepted |
| [0007](0007-better-auth.md)                 | Better Auth for accounts and sessions           | Accepted |
| [0008](0008-fumadocs-and-scalar.md)         | Fumadocs for docs, Scalar for the API reference | Accepted |
| [0009](0009-river-and-claude-candidates.md) | River for jobs; Claude candidates per concept   | Accepted |
| [0010](0010-language-detection.md)          | Our own n-gram language detector                | Accepted |
| [0011](0011-api-keys.md)                    | Our own API keys, hashed with SHA-256           | Accepted |
