import type { Schemas } from "@vernacular/api-client"

/** The message to show for a failed API call: the API's own message, or a generic one. */
export function apiErrorMessage(
  error: Schemas["Error"] | undefined,
  fallback = "The API is unreachable.",
) {
  return error?.message ?? fallback
}
