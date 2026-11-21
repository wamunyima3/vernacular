import { describe, expect, it } from "vitest"

import { groupByWave, type Language } from "./group-by-wave"

const lang = (code: string, wave: number): Language => ({
  code,
  name: code,
  own_name: code,
  wave,
  status: wave === 1 ? "active" : "planned",
  entries: { verified: 0, total: 0 },
})

describe("groupByWave", () => {
  it("orders waves ascending and keeps order inside a wave", () => {
    const result = groupByWave([lang("toi", 2), lang("bem", 1), lang("loz", 2), lang("nya", 1)])

    expect(result.map((w) => w.wave)).toEqual([1, 2])
    expect(result[0]?.languages.map((l) => l.code)).toEqual(["bem", "nya"])
    expect(result[1]?.languages.map((l) => l.code)).toEqual(["toi", "loz"])
  })

  it("returns no waves for no languages", () => {
    expect(groupByWave([])).toEqual([])
  })
})
