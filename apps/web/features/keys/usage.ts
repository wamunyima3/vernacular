import type { Schemas } from "@vernacular/api-client"

type Day = Schemas["KeyUsageDay"]

/** Today's date in UTC, as the API names usage days. */
export const utcToday = (now = new Date()) => now.toISOString().slice(0, 10)

/** Calls today and in total, from a key's daily usage. */
export function usageSummary(days: Day[], today = utcToday()) {
  const sum = (pick: (d: Day) => number) => days.reduce((n, d) => n + pick(d), 0)
  return {
    today: days.find((d) => d.day === today)?.calls ?? 0,
    calls: sum((d) => d.calls),
    errors: sum((d) => d.errors),
    limited: sum((d) => d.limited),
    busiest: Math.max(0, ...days.map((d) => d.calls)),
  }
}
