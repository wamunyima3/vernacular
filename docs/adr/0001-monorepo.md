# 0001. Monorepo with pnpm workspaces and Turborepo

- Status: Accepted
- Date: 2026-11-20

## Context
The project has a Go API, a Next.js app and an API contract they share, plus future client libraries. Keeping them in separate repos would mean releasing the contract across repos on every change.

## Decision
Use one repository. The Go module lives in `apps/api`, the web app in `apps/web`, shared TypeScript in `packages/*`, and the contract in `contracts/`. pnpm workspaces link the TS packages; Turborepo runs and caches JS tasks. A root `Makefile` is the single entry point for both languages.

## Consequences
- A change to the contract, the server and the UI lands in one commit. This fits the vertical-slice rule.
- CI runs on every change; Turbo caching keeps the JS side fast.
- Go tooling does not use Turbo; the Makefile bridges both.
