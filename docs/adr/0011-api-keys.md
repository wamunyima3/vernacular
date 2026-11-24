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
- The portal's own calls for signed-out visitors go out with its own key (`WEB_API_KEY`), both because data endpoints need one and because all those calls come from the web server's one IP address. They then share that key's rate limit, which in production needs raising above the free one (paid tiers, DEV-3).
- Rate limiting (P1-API-12) counts per key, per signed-in user, or per IP for anonymous calls. The client IP comes from `X-Forwarded-For` only through proxies listed in `TRUSTED_PROXIES`.
