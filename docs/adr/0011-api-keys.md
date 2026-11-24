# 0011. Our own API keys, hashed with SHA-256, checked against the contract

- Status: Accepted
- Date: 2026-11-24

## Context

Developers create and revoke their own keys (DEV-1, US-01), keys are hashed at rest (NFR-5), and rate limits apply per key (DEV-3). ADR-0007 left keys out of Better Auth's sessions. Better Auth has an API-key plugin, but its keys live in the web app's auth layer, while the Go API is what checks them on every call and rate-limits them.

## Decision

- Keys are a table in the Go `accounts` module (`api_keys`), managed through `GET/POST /v1/keys` and `DELETE /v1/keys/{id}` with a signed-in session. A user has at most 10 live keys; revoked keys are kept for the audit trail.
- A key is `vk_` and 52 random base32 characters (256 bits). Only its SHA-256 is stored, with its first 7 characters so the owner can tell keys apart. A slow password hash is not needed: a random 256-bit key cannot be guessed, and a fast hash keeps the per-call lookup a single index hit.
- Callers send the key in `X-API-Key`. `accounts.Middleware` resolves a key or a session to the caller; the request validator then checks the operation's `security` in the contract (`apiKey` or `session`), so which endpoints need what is written once, in `contracts/openapi.yaml`. Data endpoints accept either; account and review endpoints accept only a session, so a leaked key cannot act as its owner.
- `last_used_at` is kept to the hour, so a busy key does not write on every call.

## Consequences

- `GET /v1/languages` and `/healthz` stay open, so the portal and status checks work without a key.
- The portal's own calls to data endpoints (P1-WEB-02's "Try it") need a key or the visitor's session; it will use a server-side key of its own.
- Rate limiting (P1-API-12) can key on the caller's `KeyID`.
