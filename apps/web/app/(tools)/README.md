# (tools) route group

Signed-in screens for contributors, reviewers, senior reviewers and admins: account, contribute, review queue, disputes, sources and admin (see docs/architecture.md §3).

Every page here must:

1. Be listed in the `matcher` in `proxy.ts`, which redirects visitors without a session cookie to sign-in.
2. Call `requireSession` (or a role check built on it) inside a Suspense boundary, because the cookie check in the proxy is only a first pass.
3. Call the API with `apiAs(session.token)`, so the Go API checks the role and per-language approvals itself (ADR-0007).
