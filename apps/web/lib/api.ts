import "server-only"

import { createApiClient } from "@vernacular/api-client"

const baseUrl = process.env.API_URL ?? "http://localhost:8080"

/** Anonymous server-side API client. Together with `apiAs`, the only place the web app configures API access. */
export const api = createApiClient(baseUrl)

/** API client that acts as the signed-in user (ADR-0007). */
export function apiAs(sessionToken: string) {
  return createApiClient(baseUrl, { headers: { Authorization: `Bearer ${sessionToken}` } })
}
