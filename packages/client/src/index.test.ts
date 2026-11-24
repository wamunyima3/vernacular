import { describe, expect, it } from "vitest"

import { VernacularError, createVernacular } from "./index.js"

/** A fetch that records the request and answers with status and body. */
function fakeFetch(status: number, body: unknown) {
  const seen: Request[] = []
  const fetch = async (input: RequestInfo | URL, init?: RequestInit) => {
    seen.push(new Request(input, init))
    return new Response(JSON.stringify(body), {
      status,
      headers: { "Content-Type": "application/json" },
    })
  }
  return { fetch: fetch as typeof globalThis.fetch, seen }
}

describe("createVernacular", () => {
  it("sends the key and the query, and returns the data", async () => {
    const { fetch, seen } = fakeFetch(200, { data: [] })
    const v = createVernacular({ apiKey: "vk_test", baseUrl: "https://api.test", fetch })
    await expect(v.search({ language: "bem", q: "amenshi" })).resolves.toEqual({ data: [] })
    const req = seen[0]!
    expect(req.headers.get("X-API-Key")).toBe("vk_test")
    expect(req.url).toBe("https://api.test/v1/entries/search?language=bem&q=amenshi")
  })

  it("throws the API's error", async () => {
    const { fetch } = fakeFetch(404, {
      code: "not_found",
      message: "Entry not found",
      docs_url: "https://docs.test/errors#not_found",
    })
    const v = createVernacular({ apiKey: "k", baseUrl: "https://api.test", fetch })
    const err = await v.entry("0b0e4f4a-1f8e-4c39-8b0a-2a6f4d1c9e11").catch((e: unknown) => e)
    expect(err).toBeInstanceOf(VernacularError)
    expect(err).toMatchObject({ status: 404, code: "not_found", message: "Entry not found" })
  })
})
