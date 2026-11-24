import { describe, expect, it } from "vitest"

import { readRechecks, recheckUrl } from "./recheck-url"

const a = "00000000-0000-4000-8000-00000000000a"
const b = "00000000-0000-4000-8000-00000000000b"

describe("recheck URLs", () => {
  it("name the language and the skipped items", () => {
    expect(recheckUrl("bem")).toBe("/rechecks?language=bem")
    expect(recheckUrl("bem", [a, b])).toBe(`/rechecks?language=bem&skip=${a}%2C${b}`)
  })
  it("read back, dropping what is not an ID", () => {
    expect(readRechecks({})).toEqual({ language: undefined, skip: [] })
    expect(readRechecks({ language: "bem", skip: `${a},nope,${b}` })).toEqual({
      language: "bem",
      skip: [a, b],
    })
  })
})
