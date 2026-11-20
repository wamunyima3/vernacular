import { describe, expect, it } from "vitest"

import { decisionInput } from "./decision-input"

const a = "11111111-1111-4111-8111-111111111111"
const b = "22222222-2222-4222-8222-222222222222"

function form(entries: [string, string][]) {
  const f = new FormData()
  for (const [k, v] of entries) f.append(k, v)
  return f
}

describe("decisionInput", () => {
  it("accepts the chosen candidate only", () => {
    expect(
      decisionInput(
        form([
          ["candidate", a],
          ["candidate", b],
          ["choice", b],
          ["decision", "accept"],
        ]),
      ),
    ).toEqual({
      ok: true,
      body: { decisions: [{ sense_id: b, decision: "accept" }] },
    })
  })

  it("needs a choice from the card to accept", () => {
    expect(
      decisionInput(
        form([
          ["candidate", a],
          ["decision", "accept"],
        ]),
      ).ok,
    ).toBe(false)
    expect(
      decisionInput(
        form([
          ["candidate", a],
          ["choice", b],
          ["decision", "accept"],
        ]),
      ).ok,
    ).toBe(false)
  })

  it("rejects every candidate with the reason and note", () => {
    const f = form([
      ["candidate", a],
      ["candidate", b],
      ["decision", "reject"],
      ["reason", "Misspelled"],
      ["note", " amenshi "],
    ])
    expect(decisionInput(f)).toEqual({
      ok: true,
      body: {
        decisions: [
          { sense_id: a, decision: "reject" },
          { sense_id: b, decision: "reject" },
        ],
        comment: "Misspelled: amenshi",
      },
    })
    expect(
      decisionInput(
        form([
          ["candidate", a],
          ["decision", "reject"],
        ]),
      ).ok,
    ).toBe(false)
  })
})
