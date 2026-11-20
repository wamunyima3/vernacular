import "server-only"

import { createApiClient } from "@vernacular/api-client"

/** Server-side API client. The only place the web app configures API access. */
export const api = createApiClient(process.env.API_URL ?? process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080")
