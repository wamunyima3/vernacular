import createClient, { type Client } from "openapi-fetch"

import type { components, operations, paths } from "./schema.js"

/** Every schema in the API contract, by name: `Schemas["Entry"]`. */
export type Schemas = components["schemas"]
export type { paths }

type Query<Op extends keyof operations> = NonNullable<operations[Op]["parameters"]["query"]>
type Body<Op extends keyof operations> = operations[Op] extends {
  requestBody?: { content: { "application/json": infer B } }
}
  ? B
  : never

/** What the API said when a call failed: its error code, message and documentation link. */
export class VernacularError extends Error {
  override readonly name = "VernacularError"

  constructor(
    readonly status: number,
    readonly code: string,
    message: string,
    readonly docsUrl?: string,
  ) {
    super(message)
  }
}

export type VernacularOptions = {
  /** Your API key, from the portal's API keys page. */
  apiKey: string
  /** The API's address, such as https://api.example.org. */
  baseUrl: string
  /** A fetch to use instead of the global one, for tests or older runtimes. */
  fetch?: typeof globalThis.fetch
}

/** The result of an openapi-fetch call, or its error as a VernacularError. */
async function unwrap<T>(
  call: Promise<{ data?: T; error?: unknown; response: Response }>,
): Promise<T> {
  const { data, error, response } = await call
  if (response.ok) return data as T
  const e = (error ?? {}) as Partial<Schemas["Error"]>
  throw new VernacularError(
    response.status,
    e.code ?? "unknown",
    e.message ?? response.statusText,
    e.docs_url,
  )
}

/**
 * A client for the Vernacular API. The named methods cover the public read
 * endpoints and throw a VernacularError when a call fails; `raw` is the
 * underlying typed client for every other path.
 *
 * ```ts
 * const vernacular = createVernacular({ apiKey, baseUrl })
 * const { data } = await vernacular.search({ language: "bem", q: "amenshi" })
 * ```
 */
export function createVernacular({ apiKey, baseUrl, fetch }: VernacularOptions) {
  const raw: Client<paths> = createClient<paths>({
    baseUrl,
    fetch,
    headers: { "X-API-Key": apiKey },
  })

  return {
    raw,
    /** The languages served, with their status (GET /v1/languages). */
    languages: () => unwrap(raw.GET("/v1/languages")),
    /** Search entries by word, variant or stem (GET /v1/entries/search). */
    search: (query: Query<"searchEntries">) =>
      unwrap(raw.GET("/v1/entries/search", { params: { query } })),
    /** One entry with its senses, examples and trust fields (GET /v1/entries/{id}). */
    entry: (id: string, query?: Query<"getEntry">) =>
      unwrap(raw.GET("/v1/entries/{id}", { params: { path: { id }, query } })),
    /** A concept and its words across languages (GET /v1/concepts/{id}). */
    concept: (id: string, query?: Query<"getConcept">) =>
      unwrap(raw.GET("/v1/concepts/{id}", { params: { path: { id }, query } })),
    /** Words for English glosses, in one or more languages (GET /v1/lookup). */
    lookup: (query: Query<"lookup">) => unwrap(raw.GET("/v1/lookup", { params: { query } })),
    /** A phrase pack in one language (GET /v1/phrases). */
    phrases: (query: Query<"getPhrasePack">) =>
      unwrap(raw.GET("/v1/phrases", { params: { query } })),
    /** Verified phrase translations (POST /v1/translate). */
    translate: (body: Body<"translate">) => unwrap(raw.POST("/v1/translate", { body })),
    /** Which language a text is in (POST /v1/detect). */
    detect: (body: Body<"detectLanguage">) => unwrap(raw.POST("/v1/detect", { body })),
    /** Published dataset releases (GET /v1/releases). */
    releases: () => unwrap(raw.GET("/v1/releases")),
  }
}

export type Vernacular = ReturnType<typeof createVernacular>
