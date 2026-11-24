import { describe, expect, it } from "vitest"

import { usageSummary, utcToday } from "./usage"

describe("usageSummary", () => {
  it("adds up the days and finds today", () => {
    const days = [
      { day: "2026-11-24", calls: 12, errors: 2, limited: 1 },
      { day: "2026-11-20", calls: 30, errors: 0, limited: 0 },
    ]
    expect(usageSummary(days, "2026-11-24")).toEqual({
      today: 12,
      calls: 42,
      errors: 2,
      limited: 1,
      busiest: 30,
    })
    expect(usageSummary([], "2026-11-24")).toMatchObject({ today: 0, calls: 0, busiest: 0 })
  })

  it("names today in UTC", () => {
    expect(utcToday(new Date("2026-11-24T23:30:00-02:00"))).toBe("2026-11-25")
  })
})
