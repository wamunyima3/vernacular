import type { Schemas } from "@vernacular/api-client"

export const reviewerLevels = [
  "reviewer",
  "senior",
] as const satisfies readonly Schemas["ReviewerLevel"][]
