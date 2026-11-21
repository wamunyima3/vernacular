import type { Schemas } from "@vernacular/api-client"

import { formText, isOneOf } from "./form-data"

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

/**
 * Reads the fields of WordClassFields: the part of speech, and the noun class
 * for nouns only. Empty or unknown values are left out.
 */
export function wordClassInput(
  form: FormData,
  prefix = "",
): { part_of_speech?: Schemas["PartOfSpeech"]; noun_class?: string } {
  const text = formText(form)
  const pos = text(`${prefix}part_of_speech`)
  if (!isOneOf(partsOfSpeech, pos)) return {}
  const nounClass = pos === "noun" ? text(`${prefix}noun_class`) : ""
  return { part_of_speech: pos, ...(nounClass && { noun_class: nounClass }) }
}

/** "noun · class 7/8", or "" when neither is known. */
export function wordClassLabel(w: { part_of_speech?: string; noun_class?: string }): string {
  return [w.part_of_speech, w.noun_class && `class ${w.noun_class}`].filter(Boolean).join(" · ")
}
