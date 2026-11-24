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

  it("gives an example's English and the word it uses", () => {
    expect(
      candidateDetail({
        item_id: "a",
        headword: "Nalanwa amenshi",
        translation_en: "I drank water",
        word: "amenshi",
      }),
    ).toBe("“I drank water” · uses amenshi")
  })

  it("is empty for a bare candidate", () => {
    expect(candidateDetail({ item_id: "a", headword: "x" })).toBe("")
  })
})
