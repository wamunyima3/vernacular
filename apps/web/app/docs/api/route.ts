import { readFileSync } from "node:fs"
import { join } from "node:path"

import { ApiReference } from "@scalar/nextjs-api-reference"

import { baseUrl } from "@/lib/api"

// The contract itself, so the reference cannot drift from it (ADR-0003, ADR-0008). The
// route is prerendered as a static page by `next build`; Next runs from apps/web.
const contract = readFileSync(join(process.cwd(), "../../contracts/openapi.yaml"), "utf8")

// Where browsers send the request runner's calls; the server's own address by default.
const apiUrl = process.env.PUBLIC_API_URL ?? baseUrl

/** The interactive API reference (DEV-2, US-02): every endpoint, with a request runner. */
export const GET = ApiReference({
  content: contract,
  servers: [{ url: apiUrl, description: "Vernacular API" }],
  authentication: { preferredSecurityScheme: "apiKey" },
  metaData: { title: "API reference · Vernacular" },
  // No vendor toolbar or hosted AI chat on a public page.
  showDeveloperTools: "never",
  agent: { disabled: true },
})
