/**
 * Global account roles (ADR-0007). Reviewer and senior reviewer are per-language
 * approvals, not roles, so they are not listed here.
 */
export const roles = {
  developer: "developer",
  contributor: "contributor",
  admin: "admin",
} as const

export type Role = (typeof roles)[keyof typeof roles]
