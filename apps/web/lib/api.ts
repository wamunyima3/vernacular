import "server-only"

import { createApiClient } from "@vernacular/api-client"

/** The API as the web server reaches it. */
export const baseUrl = process.env.API_URL ?? "http://localhost:8080"

const portalKey = process.env.WEB_API_KEY

/**
 * Server-side API client for pages anyone can see. It calls with the portal's own API
 * key when WEB_API_KEY is set, so visitors do not share one IP's rate limit (ADR-0011).
 * Together with `apiAs`, the only place the web app configures API access.
 */
export const api = createApiClient(
  baseUrl,
  portalKey ? { headers: { "X-API-Key": portalKey } } : undefined,
)

/** API client that acts as the signed-in user (ADR-0007). */
export function apiAs(sessionToken: string) {
  return createApiClient(baseUrl, { headers: { Authorization: `Bearer ${sessionToken}` } })
}
