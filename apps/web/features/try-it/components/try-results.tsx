import type { Schemas } from "@vernacular/api-client"
import Link from "next/link"

import { ErrorCard } from "@/components/error-card"
import { NoticeCard } from "@/components/notice-card"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { apiErrorMessage } from "@/lib/api-error"
import { wordClassLabel } from "@/lib/word-class"

import { apiForVisitor } from "../../auth/session"
import { TrustBadge } from "../../trust/components/trust-badge"
import { JsonPanel } from "./json-panel"
import { searchParamsOf, type TryQuery } from "./try-form"

const matchLabels: Record<Schemas["Match"], string> = {
  headword: "Matched the headword",
  variant: "Matched a spelling variant",
  root: "Matched the root",
  similar: "Similar spelling",
}

/**
 * Calls GET /v1/entries/search as the visitor and shows the result two ways: formatted,
 * and as the raw JSON. Side by side on wide screens, as tabs on phones.
 */
export async function TryResults({
  query,
  languageName,
}: {
  query: TryQuery
  languageName: string
}) {
  const api = await apiForVisitor()
  const { data, error, response } = await api.GET("/v1/entries/search", {
    params: { query: { language: query.language, q: query.q, min_status: query.minStatus } },
  })
  if (response.status === 401) {
    return (
      <NoticeCard
        title="Sign in to try the API"
        actions={
          <Button nativeButton={false} render={<Link href="/sign-in?next=/try" />}>
            Sign in
          </Button>
        }
      >
        Searching here calls the live API, which needs an account or an API key.
      </NoticeCard>
    )
  }
  if (!data) return <ErrorCard title="Search failed" message={apiErrorMessage(error)} />

  const request = `GET /v1/entries/search?${searchParamsOf(query)}`
  const json = <JsonPanel json={JSON.stringify(data, null, 2)} request={request} />
  const result =
    data.data.length === 0 ? (
      <NoResults query={query} languageName={languageName} />
    ) : (
      <ol className="flex flex-col gap-3">
        {data.data.map((hit) => (
          <Hit key={hit.entry.id} hit={hit} />
        ))}
      </ol>
    )

  return (
    <>
      <div className="hidden gap-6 lg:grid lg:grid-cols-2">
        <section aria-label="Result">{result}</section>
        {json}
      </div>
      <Tabs defaultValue="result" className="lg:hidden">
        <TabsList>
          <TabsTrigger value="result">Result</TabsTrigger>
          <TabsTrigger value="json">JSON</TabsTrigger>
        </TabsList>
        <TabsContent value="result">{result}</TabsContent>
        <TabsContent value="json">{json}</TabsContent>
      </Tabs>
    </>
  )
}

function Hit({ hit: { match, entry } }: { hit: Schemas["EntryHit"] }) {
  const wordClass = wordClassLabel(entry)
  return (
    <li className="flex flex-col gap-2 rounded-lg border p-4">
      <div className="flex flex-wrap items-baseline gap-x-2">
        <span className="text-xl font-semibold" lang={entry.language}>
          {entry.headword}
        </span>
        {wordClass && <span className="text-sm text-muted-foreground">{wordClass}</span>}
      </div>
      <ul className="flex flex-col gap-1">
        {entry.senses.map((s) => (
          <li key={s.id} className="flex flex-wrap items-center gap-2 text-pretty">
            {s.gloss} <TrustBadge status={s.trust.status} />
          </li>
        ))}
      </ul>
      <p className="text-xs text-muted-foreground">{matchLabels[match]}</p>
    </li>
  )
}

/** docs/design/ui.md: one action, to widen the search, or to add the word. */
function NoResults({ query, languageName }: { query: TryQuery; languageName: string }) {
  const strict = query.minStatus === "verified"
  const wider = `/try?${searchParamsOf({ ...query, minStatus: "machine" })}`
  return (
    <NoticeCard
      title={`No ${strict ? "verified " : ""}entry for “${query.q}” in ${languageName}.`}
      actions={
        strict ? (
          <Button nativeButton={false} render={<Link href={wider} />}>
            Include unverified results
          </Button>
        ) : (
          <Button nativeButton={false} render={<Link href="/contribute" />}>
            Contribute this word
          </Button>
        )
      }
    >
      {strict
        ? "Unverified words may exist, from sources or machine translation."
        : "Nothing matches, even among unchecked words."}
    </NoticeCard>
  )
}

export function TryResultsSkeleton() {
  return <Skeleton className="h-48 w-full" />
}
