import { describe, expect, it } from "vitest"

import { contributionInput } from "./form-data"

function form(entries: Record<string, string>) {
  const f = new FormData()
  for (const [k, v] of Object.entries(entries)) f.append(k, v)
  return f
}

const base = {
  language: "bem",
  headword: " icipuna ",
  part_of_speech: "noun",
  noun_class: "7/8",
  gloss_en: "chair",
  agree_to_licence: "on",
}

describe("contributionInput", () => {
  it("passes on the meaning a correction names", () => {
    const id = "0b0e4f4a-1f8e-4c39-8b0a-2a6f4d1c9e11"
    const r = contributionInput(form({ ...base, corrects: id }))
    expect(r.ok && r.body.corrects).toBe(id)
  })

  it("builds the contract input, without empty optional fields", () => {
    expect(contributionInput(form(base))).toEqual({
      ok: true,
      body: {
        language: "bem",
        headword: "icipuna",
        part_of_speech: "noun",
        noun_class: "7/8",
        gloss_en: "chair",
        agree_to_licence: true,
      },
    })
  })

  it("drops the noun class for other parts of speech", () => {
    const r = contributionInput(form({ ...base, part_of_speech: "verb" }))
    expect(r.ok && "noun_class" in r.body).toBe(false)
  })

  it("needs both halves of an example, and agreement to the licence", () => {
    const unagreed = { ...base, agree_to_licence: "", example_text: "Ikala pa cipuna." }
    expect(contributionInput(form(unagreed))).toEqual({
      ok: false,
      errors: {
        example: expect.stringContaining("translation"),
        agree_to_licence: expect.stringContaining("licence"),
      },
    })
  })
})
