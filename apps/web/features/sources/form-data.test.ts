import { describe, expect, it } from "vitest"

import { sourceInput } from "./form-data"

function form(entries: [string, string][]) {
  const f = new FormData()
  for (const [k, v] of entries) f.append(k, v)
  return f
}

describe("sourceInput", () => {
  it("trims text, keeps known uses, omits an empty url", () => {
    const input = sourceInput(
      form([
        ["kind", "licensed"],
        ["name", "  GATITOS "],
        ["rights_holder", "Google"],
        ["licence", "CC BY 4.0"],
        ["url", ""],
        ["allowed_uses", "derive"],
        ["allowed_uses", "sell-forever"],
      ]),
    )
    expect(input).toEqual({
      kind: "licensed",
      name: "GATITOS",
      rights_holder: "Google",
      licence: "CC BY 4.0",
      allowed_uses: ["derive"],
      notes: "",
    })
  })

  it("falls back to licensed for an unknown kind", () => {
    expect(sourceInput(form([["kind", "stolen"]])).kind).toBe("licensed")
  })
})
