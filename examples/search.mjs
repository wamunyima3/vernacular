// Search the Bemba dictionary for "amenshi", unchecked words included. Node 18 or later.
const url = new URL("/v1/entries/search", process.env.VERNACULAR_URL)
url.search = new URLSearchParams({ language: "bem", q: "amenshi", min_status: "machine" })

const res = await fetch(url, { headers: { "X-API-Key": process.env.VERNACULAR_KEY } })
const body = await res.json()
if (!res.ok) throw new Error(`${res.status} ${body.code}: ${body.message}`)

for (const { entry } of body.data) {
  const meanings = entry.senses.map((s) => `${s.gloss} (${s.trust.status})`)
  console.log(`${entry.headword}: ${meanings.join(", ")}`)
}
