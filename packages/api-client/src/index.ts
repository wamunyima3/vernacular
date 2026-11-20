import createClient, { type ClientOptions } from "openapi-fetch"

import type { components, paths } from "./schema"

/** Every schema in the contract, by name: `Schemas["Language"]`. */
export type Schemas = components["schemas"]

/** Creates a typed client. Paths, parameters and responses all come from the contract. */
export function createApiClient(baseUrl: string, options?: Omit<ClientOptions, "baseUrl">) {
  return createClient<paths>({ ...options, baseUrl })
}
