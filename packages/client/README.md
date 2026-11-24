# @vernacular/client

JavaScript and TypeScript client for the Vernacular API: human-verified words and phrases in
Zambian languages, each with its trust status and sources.

```bash
npm install @vernacular/client
```

```ts
import { createVernacular, VernacularError } from "@vernacular/client"

const vernacular = createVernacular({
  apiKey: process.env.VERNACULAR_KEY!,
  baseUrl: process.env.VERNACULAR_URL!,
})

const { data } = await vernacular.search({ language: "bem", q: "amenshi" })
const pack = await vernacular.phrases({ pack: "greetings", lang: "nya" })

try {
  await vernacular.entry("…")
} catch (e) {
  if (e instanceof VernacularError) console.log(e.status, e.code, e.docsUrl)
}
```

Methods: `languages`, `search`, `entry`, `concept`, `lookup`, `phrases`, `translate`, `detect` and
`releases`. Each returns the response body, typed from the API contract, and throws a
`VernacularError` when the call fails. Every other endpoint is on `vernacular.raw`, an
[openapi-fetch](https://openapi-ts.dev/openapi-fetch/) client with the key already set:

```ts
const { data, error } = await vernacular.raw.GET("/v1/phrases/export", {
  params: { query: { pack: "interface", lang: "bem", format: "po" } },
  parseAs: "text",
})
```

Only `verified` items come back unless you pass `min_status`. Read the trust fields before you
rely on a word.

Licence: MIT. The data you receive has its own licence; see the API documentation.
