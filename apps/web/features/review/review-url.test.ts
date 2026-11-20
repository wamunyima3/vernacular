import { describe, expect, it } from "vitest"

import { nextSession, readSession, reviewUrl } from "./review-url"

const a = "11111111-1111-4111-8111-111111111111"
const b = "22222222-2222-4222-8222-222222222222"

describe("review session URL", () => {
  it("round-trips and drops malformed values", () => {
    const url = reviewUrl({ language: "bem", skip: [a, b], done: 3 })
    const params = Object.fromEntries(new URL(url, "http://x").searchParams)
    expect(readSession(params)).toEqual({ language: "bem", skip: [a, b], done: 3 })
    expect(readSession({ language: "bem", skip: "nope", done: "-2" })).toEqual({
      language: "bem",
      skip: [],
      done: 0,
    })
  })

  it("skips the card either way and counts only decisions", () => {
    const s = { language: "bem", skip: [a], done: 1 }
    expect(nextSession(s, b, true)).toEqual({ language: "bem", skip: [a, b], done: 2 })
    expect(nextSession(s, b, false)).toEqual({ language: "bem", skip: [a, b], done: 1 })
  })
})
