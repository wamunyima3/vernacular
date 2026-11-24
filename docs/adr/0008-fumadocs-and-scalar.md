# ADR-0008: Fumadocs for docs, Scalar for the API reference

- **Status:** Accepted
- **Date:** 2026-11-20
- **Task:** P0-ENG-09

## Context

The portal needs written docs (getting started, verification levels, concepts, varieties, licences) and an interactive API reference generated from `contracts/openapi.yaml` (DEV-2, US-02). Both should live inside `apps/web`, so that there is one deployable and one design system.

| Criterion                 | Fumadocs 16                       | Nextra 4                                |
| ------------------------- | --------------------------------- | --------------------------------------- |
| Licence                   | MIT                               | MIT                                     |
| Last release (2026-11-20) | 2026-11-19, weekly                | 2025-12-04, about a year ago            |
| Next.js 16, App Router    | Peer dependency `next 16.x`       | `>=14`, not tested against 16 features  |
| Inside an existing app    | Yes: a library mounted at `/docs` | Takes over the app's routing and layout |
| Styling                   | Tailwind v4, shadcn-compatible    | Its own theme                           |

For the API reference, Scalar (MIT, released daily) renders an OpenAPI file with a built-in request runner and code samples in many languages. Fumadocs' own OpenAPI plugin renders static pages, but it has no request runner.

## Decision

- **Fumadocs** (`fumadocs-core`, `fumadocs-ui`, `fumadocs-mdx`) for MDX docs under `apps/web/content/docs`, served at `/docs`.
- **Scalar** (`@scalar/nextjs-api-reference`) at `/docs/api`. It reads `contracts/openapi.yaml` at build time, so the reference never drifts from the contract.
- `/docs/errors` stays a typed page, because it is checked against `ErrorCode` at compile time.

### Amendment, 2026-11-24 (P1-WEB-04)

The docs use `fumadocs-mdx` and `fumadocs-core` only: MDX compiled at build time, and a page tree whose order `content/docs/meta.json` sets. `fumadocs-ui` is not used, because its layout brings a second header, sidebar and stylesheet into the portal. The pages render with the app's own components (`features/docs/mdx-components.tsx`), keeping one design system. Fumadocs UI's search and table of contents are given up until the docs are long enough to need them.

## Consequences

- Docs are versioned with the code and reviewed in the same pull requests.
- Two documentation libraries, each doing one job. If Fumadocs' OpenAPI plugin gains a request runner, revisit this decision and drop Scalar.
