import type { Schemas } from "@vernacular/api-client"

import { formText } from "@/lib/form-data"
import { wordClassInput } from "@/lib/word-class"

export const licenceError = "Agree to publish your contribution under the dataset licence."

/** The example fields of a form, trimmed; either half may be empty. */
export function exampleInput(form: FormData): Schemas["ExampleInput"] {
  const text = formText(form)
  return { text: text("example_text"), translation_en: text("example_translation") }
}

export type ContributionErrors = Partial<Record<"example" | "agree_to_licence", string>>

/**
 * Reads the contribute form into the contract's input, or field errors the
 * contract cannot express. The API validates the rest.
 */
export function contributionInput(
  form: FormData,
): { ok: true; body: Schemas["ContributionInput"] } | { ok: false; errors: ContributionErrors } {
  const text = formText(form)
  const errors: ContributionErrors = {}
  const { part_of_speech = "noun", noun_class } = wordClassInput(form)
  const example = exampleInput(form)
  if (Boolean(example.text) !== Boolean(example.translation_en)) {
    errors.example = "Give the example and its English translation, or leave both empty."
  }
  if (form.get("agree_to_licence") !== "on") {
    errors.agree_to_licence = licenceError
  }
  if (Object.keys(errors).length > 0) return { ok: false, errors }
  const corrects = text("corrects")

  return {
    ok: true,
    body: {
      language: text("language"),
      headword: text("headword"),
      part_of_speech,
      gloss_en: text("gloss_en"),
      ...(noun_class && { noun_class }),
      ...(example.text && { example }),
      ...(corrects && { corrects }),
      agree_to_licence: true,
    },
  }
}
