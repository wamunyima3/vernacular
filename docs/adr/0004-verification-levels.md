# 0004. Verification levels and shared trust fields

- Status: Accepted
- Date: 2026-11-20

## Context

The product rule is that every item says where it came from and who checked it (§1, PRV-1, §7). Model output must never be passed off as verified (DQ-2).

## Decision

- Five statuses: `verified`, `reviewed`, `sourced`, `machine`, `disputed`.
- One `Trust` component in the contract (`status`, `reviews`, `sources[]` with licence, `release`) is embedded in every item schema.
- Read endpoints filter with `min_status`, defaulting to `verified`. `machine` is returned only on request, and `disputed` is never returned.
- Exactly one function in the `review` module may promote an item to `verified`. It requires two accepts from different, approved reviewers, neither of whom authored the item.

## Consequences

- Developers learn the trust fields once.
- The promotion function is the most heavily tested code in the repo.
- Adding a status is a breaking contract change.

## Amendment (2026-11-21): `rejected`, and how statuses are derived

Nothing said what happens to an item two reviewers reject. Returning it to `sourced` would keep it in the queue forever and still serve it on request. The database gains a sixth status, `rejected`. Rejected items are never returned, so the contract keeps its five statuses and this is not a breaking change. They are kept for history and for acceptance rates per provider (AI-5).

The promotion function derives an item's status from its origin and its reviews, so recording a decision and undoing one go through the same rule:

| Reviews (distinct, approved, not the author) | Status                                                            |
| -------------------------------------------- | ----------------------------------------------------------------- |
| none, or one reject                          | origin: `machine` if only machine sources cite it, else `sourced` |
| one accept                                   | `reviewed`                                                        |
| two or more accepts, no rejects              | `verified`                                                        |
| two or more rejects, no accepts              | `rejected`                                                        |
| accepts and rejects                          | `disputed`, for a senior reviewer (P1-REV-04)                     |
