import type { Schemas } from "@vernacular/api-client"

export type CardKind = Schemas["ReviewCard"]["kind"]

/** How a review card of each kind names its candidates, its English and a correction. */
export const kinds: Record<
  CardKind,
  {
    noun: string
    heading: string
    wraps: boolean // candidates are sentences or phrases, not single words
    label: string // the "Write my own" field
    maxLength: number
    rejectWhy: (gloss: string) => string
  }
> = {
  sense: {
    noun: "word",
    heading: "English",
    wraps: false,
    label: "The right word",
    maxLength: 100,
    rejectWhy: (g) => `None of them means “${g}”. Why?`,
  },
  phrase: {
    noun: "phrase",
    heading: "English phrase",
    wraps: true,
    label: "The right phrase",
    maxLength: 200,
    rejectWhy: (g) => `None of them means “${g}”. Why?`,
  },
  example: {
    noun: "sentence",
    heading: "Example sentences for",
    wraps: true,
    label: "Your own sentence using the word",
    maxLength: 200,
    rejectWhy: (g) => `None of these sentences is a good example of “${g}”. Why?`,
  },
}

export const isCardKind = (v: unknown): v is CardKind => typeof v === "string" && v in kinds
