import type { Schemas } from "@vernacular/api-client"

import { isOneOf } from "@/lib/form-data"

/** Where the reviewer is in a session: the language, skipped concepts and how many cards they have done. */
export type ReviewSession = { language: string; skip: string[]; done: number }

/** The reviews just recorded, which can still be undone. */
export type LastDecision = { reviews: string[]; decision: Schemas["ReviewDecision"] }

/** Search params as Next.js passes them. */
export type Params = Record<string, string | string[] | undefined>

const decisions = [
  "accept",
  "reject",
  "correct",
] as const satisfies readonly Schemas["ReviewDecision"][]
const uuid = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i

const one = (params: Params, k: string) => (Array.isArray(params[k]) ? params[k][0] : params[k])
const ids = (params: Params, k: string) =>
  (one(params, k) ?? "").split(",").filter((s) => uuid.test(s))

/** Reads the session from URL or form values, dropping anything malformed. */
export function readSession(
  params: Params,
): Omit<ReviewSession, "language"> & { language?: string } {
  const done = Number(one(params, "done"))
  return {
    language: one(params, "language"),
    skip: ids(params, "skip").slice(-100),
    done: Number.isInteger(done) && done > 0 ? done : 0,
  }
}

/** Reads the decision that can still be undone, if any. */
export function readLast(params: Params): LastDecision | undefined {
  const decision = one(params, "decided")
  const reviews = ids(params, "undo").slice(0, 50)
  return isOneOf(decisions, decision) && reviews.length > 0 ? { reviews, decision } : undefined
}

/**
 * The session after leaving a card. The concept is skipped for the rest of the
 * session either way: candidates the reviewer did not choose stay in the queue
 * for others instead of being rejected.
 */
export function nextSession(s: ReviewSession, conceptId: string, decided: boolean): ReviewSession {
  return { ...s, skip: [...s.skip, conceptId], done: s.done + (decided ? 1 : 0) }
}

/** The session before the last decision: its card back, and one fewer done. */
export function undoneSession(s: ReviewSession): ReviewSession {
  return { ...s, skip: s.skip.slice(0, -1), done: Math.max(0, s.done - 1) }
}

/** The URL of a session; empty parts are left out. */
export function reviewUrl({ language, skip, done }: ReviewSession, last?: LastDecision): string {
  const q = new URLSearchParams({ language })
  if (skip.length > 0) q.set("skip", skip.join(","))
  if (done > 0) q.set("done", String(done))
  if (last) {
    q.set("decided", last.decision)
    q.set("undo", last.reviews.join(","))
  }
  return `/review?${q}`
}
