import { describe, expect, it } from "vitest"

import { candidateDetail } from "./word-class"

describe("candidateDetail", () => {
  it("names the meaning a correction replaces", () => {
    expect(
      candidateDetail({
        item_id: "a",
        headword: "amensi",
        part_of_speech: "noun",
        noun_class: "6",
        corrects: { item_id: "b", headword: "amenshi", gloss: "water" },
      }),
    ).toMatch(/· replaces amenshi \(water\)$/)
  })

  it("is empty for a bare candidate", () => {
    expect(candidateDetail({ item_id: "a", headword: "x" })).toBe("")
  })
})
