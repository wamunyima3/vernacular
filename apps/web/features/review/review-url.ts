/** Where the reviewer is in a session: the language, skipped concepts and how many cards they have done. */
export type ReviewSession = { language: string; skip: string[]; done: number }

const uuid = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i

/** Reads the session from URL or form values, dropping anything malformed. */
export function readSession(
  params: Record<string, string | string[] | undefined>,
): Omit<ReviewSession, "language"> & { language?: string } {
  const one = (k: string) => (Array.isArray(params[k]) ? params[k][0] : params[k])
  const done = Number(one("done"))
  return {
    language: one("language"),
    skip: (one("skip") ?? "")
      .split(",")
      .filter((s) => uuid.test(s))
      .slice(-100),
    done: Number.isInteger(done) && done > 0 ? done : 0,
  }
}

/**
 * The session after leaving a card. The concept is skipped for the rest of the
 * session either way: candidates the reviewer did not choose stay in the queue
 * for others instead of being rejected.
 */
export function nextSession(s: ReviewSession, conceptId: string, decided: boolean): ReviewSession {
  return { ...s, skip: [...s.skip, conceptId], done: s.done + (decided ? 1 : 0) }
}

/** The URL of a session; skip and done are left out when empty. */
export function reviewUrl({ language, skip, done }: ReviewSession): string {
  const q = new URLSearchParams({ language })
  if (skip.length > 0) q.set("skip", skip.join(","))
  if (done > 0) q.set("done", String(done))
  return `/review?${q}`
}
