import createClient from "openapi-fetch"

import type { components, paths } from "./schema"

export type { components, paths }

/** Every schema in the contract, by name: `Schemas["Language"]`. */
export type Schemas = components["schemas"]

export type ApiClient = ReturnType<typeof createApiClient>

/** Creates a typed client. Paths, parameters and responses all come from the contract. */
export function createApiClient(baseUrl: string, init?: { headers?: HeadersInit }) {
  return createClient<paths>({ baseUrl, headers: init?.headers })
}
