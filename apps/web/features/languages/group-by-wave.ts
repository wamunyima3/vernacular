import type { Schemas } from "@vernacular/api-client"

export type Language = Schemas["Language"]

export type Wave = { wave: number; languages: Language[] }

/** Groups languages by rollout wave, in wave order, keeping each wave's input order. */
export function groupByWave(languages: readonly Language[]): Wave[] {
  const byWave = Map.groupBy(languages, (l) => l.wave)
  return [...byWave.entries()]
    .sort(([a], [b]) => a - b)
    .map(([wave, langs]) => ({ wave, languages: langs }))
}
