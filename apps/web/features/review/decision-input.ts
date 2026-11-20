import type { Schemas } from "@vernacular/api-client"

import { formText, isOneOf } from "@/lib/form-data"

export const rejectReasons = [
  "Wrong meaning",
  "Misspelled",
  "Not used in this language",
  "Offensive",
  "Other",
] as const

/**
 * Reads a decision form into the contract's input. Accept takes the chosen
 * candidate; reject applies to every candidate on the card.
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
        decisions: candidates.map((sense_id) => ({ sense_id, decision: "reject" as const })),
        comment: note ? `${reason}: ${note}` : reason,
      },
    }
  }
  const choice = text("choice")
  if (!candidates.includes(choice))
    return { ok: false, error: "Choose the right word, or reject them all." }
  return { ok: true, body: { decisions: [{ sense_id: choice, decision: "accept" }] } }
}
