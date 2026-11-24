import type { Schemas } from "@vernacular/api-client"

import { isOneOf } from "./form-data"

export type MinStatus = Schemas["MinStatus"]

/** The verification levels a caller can ask for, strictest first (ADR-0004). */
export const minStatuses: { value: MinStatus; label: string }[] = [
  { value: "verified", label: "Verified only" },
  { value: "reviewed", label: "Reviewed once or more" },
  { value: "sourced", label: "Include unchecked sources" },
  { value: "machine", label: "Include machine output" },
]

const values = minStatuses.map((s) => s.value)

/** A min_status from the URL, or the API's default, verified. */
export const parseMinStatus = (v: unknown): MinStatus => (isOneOf(values, v) ? v : "verified")
