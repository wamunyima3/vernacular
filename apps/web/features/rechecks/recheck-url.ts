import { idList, param, type Params } from "@/lib/search-params"

/** The re-check page for a language, with the items skipped so far. */
export function recheckUrl(language: string, skip: string[] = []): string {
  const q = new URLSearchParams({ language })
  if (skip.length > 0) q.set("skip", skip.join(","))
  return `/rechecks?${q}`
}

/** The language asked for and the last 100 items skipped, dropping anything malformed. */
export function readRechecks(params: Params): { language?: string; skip: string[] } {
  return { language: param(params, "language"), skip: idList(params, "skip").slice(-100) }
}
