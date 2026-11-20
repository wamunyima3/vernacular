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
