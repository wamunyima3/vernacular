import type { Schemas } from "@vernacular/api-client"

import { formText, isOneOf } from "@/lib/form-data"
import { wordClassInput } from "@/lib/word-class"

import { isCardKind, kinds } from "./kinds"

export const rejectReasons = [
  "Wrong meaning",
  "Misspelled",
  "Not used in this language",
  "Offensive",
  "Other",
] as const

/**
 * Reads a decision form into the contract's input. Accept takes the chosen
 * candidate; reject and correct apply to every candidate on the card.
 */
export function decisionInput(
  form: FormData,
): { ok: true; body: Schemas["ReviewInput"] } | { ok: false; error: string } {
  const text = formText(form)
  const candidates = form.getAll("candidate").map(String)
  if (form.get("decision") === "reject") {
    const reason = form.get("reason")
    if (!isOneOf(rejectReasons, reason)) return { ok: false, error: "Choose a reason." }
    const note = text("note")
    return {
      ok: true,
      body: {
        decisions: candidates.map((item_id) => ({ item_id, decision: "reject" as const })),
        comment: note ? `${reason}: ${note}` : reason,
      },
    }
  }
  if (form.get("decision") === "correct") {
    const sent = form.get("kind")
    const kind = isCardKind(sent) ? sent : "sense"
    const headword = text("correction_headword")
    if (!headword) return { ok: false, error: `Write the right ${kinds[kind].noun}.` }
    const translation_en = text("correction_translation_en")
    if (kind === "example" && !translation_en)
      return { ok: false, error: "Write the sentence's English." }
    return {
      ok: true,
      body: {
        decisions: candidates.map((item_id) => ({ item_id, decision: "correct" as const })),
        correction:
          kind === "example"
            ? { headword, translation_en }
            : kind === "sense"
              ? { headword, ...wordClassInput(form, "correction_") }
              : { headword },
      },
    }
  }
  const choice = text("choice")
  if (!candidates.includes(choice))
    return { ok: false, error: "Choose the right word, or reject them all." }
  return { ok: true, body: { decisions: [{ item_id: choice, decision: "accept" }] } }
}
