# 0003. OpenAPI-first contract with generated code

- Status: Accepted
- Date: 2026-11-20

## Context

DEV-2 needs an accurate OpenAPI description and interactive docs. The web app and future SDKs (DEV-4) need types that match the server exactly.

## Decision

`contracts/openapi.yaml` is written by hand and is the source of truth. `oapi-codegen` generates a strict chi server interface and types into `apps/api/internal/gen`. `openapi-typescript` generates types into `packages/api-client`. Nobody edits generated files. CI runs `make gen` and fails if anything changes.

## Consequences

- The server cannot drift from the docs: an unimplemented operation fails to compile.
- Every API change starts with a contract edit, which makes it easy to review.
- Shared schemas (trust fields, errors, pagination) are written once as components and reused through `$ref`.
