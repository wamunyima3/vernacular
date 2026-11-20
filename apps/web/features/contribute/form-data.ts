import type { Schemas } from "@vernacular/api-client"

import { formText, isOneOf } from "@/lib/form-data"

// A key per enum value, so a part of speech added to the contract fails the type check here.
const partOfSpeechKeys: Record<Schemas["PartOfSpeech"], null> = {
  noun: null,
  verb: null,
  adjective: null,
  adverb: null,
  pronoun: null,
  numeral: null,
  conjunction: null,
  preposition: null,
  interjection: null,
  phrase: null,
  ideophone: null,
}
export const partsOfSpeech = Object.keys(partOfSpeechKeys) as Schemas["PartOfSpeech"][]

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
  const posValue = text("part_of_speech")
  const pos = isOneOf(partsOfSpeech, posValue) ? posValue : "noun"
  const example = { text: text("example_text"), translation_en: text("example_translation") }
  if (Boolean(example.text) !== Boolean(example.translation_en)) {
    errors.example = "Give the example and its English translation, or leave both empty."
  }
  if (form.get("agree_to_licence") !== "on") {
    errors.agree_to_licence = "Agree to publish your contribution under the dataset licence."
  }
  if (Object.keys(errors).length > 0) return { ok: false, errors }

  const nounClass = pos === "noun" ? text("noun_class") : ""
  return {
    ok: true,
    body: {
      language: text("language"),
      headword: text("headword"),
      part_of_speech: pos,
      gloss_en: text("gloss_en"),
      ...(nounClass && { noun_class: nounClass }),
      ...(example.text && { example }),
      agree_to_licence: true,
    },
  }
}
