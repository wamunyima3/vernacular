/** A meaning being corrected, and what the contribute form starts with (US-18). */
export type Correcting = { corrects: string; language: string; headword: string; gloss_en: string }

const uuidPattern = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i

/** The contribute form, prefilled with a meaning and linked to it. */
export const correctionHref = (c: Correcting) => `/contribute?${new URLSearchParams(c)}`

/** The meaning a correctionHref link names, if the link is whole. */
export function correctingFrom(
  q: Record<string, string | string[] | undefined>,
): Correcting | undefined {
  const { corrects, language, headword, gloss_en } = q
  if (typeof corrects !== "string" || !uuidPattern.test(corrects)) return undefined
  if (typeof language !== "string" || typeof headword !== "string" || typeof gloss_en !== "string")
    return undefined
  return { corrects, language, headword, gloss_en }
}
