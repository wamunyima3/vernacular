import type { Metadata } from "next"
import { Suspense } from "react"

import { Skeleton } from "@/components/ui/skeleton"
import { parseMinStatus } from "@/lib/min-status"
import { ErrorCard } from "@/components/error-card"
import { getLanguages, isActive } from "@/features/languages/get-languages"
import { TryForm, type TryQuery } from "@/features/try-it/components/try-form"
import { TryResults, TryResultsSkeleton } from "@/features/try-it/components/try-results"

export const metadata: Metadata = { title: "Try it" }

const one = (v: string | string[] | undefined) => (typeof v === "string" ? v.trim() : "")

/** The form, filled from the URL, and the results of the search it describes. */
async function Search({ searchParams }: PageProps<"/try">) {
  const [params, result] = await Promise.all([searchParams, getLanguages()])
  if (!result.ok)
    return <ErrorCard title="Languages could not be loaded" message={result.message} />
  const languages = result.languages.filter(isActive)
  const language = one(params.language)
  const query: TryQuery = {
    language: languages.some((l) => l.code === language) ? language : (languages[0]?.code ?? ""),
    q: one(params.q).slice(0, 100),
    minStatus: parseMinStatus(params.min_status),
  }
  const key = JSON.stringify(query)
  const name = languages.find((l) => l.code === query.language)?.name ?? query.language
  return (
    <>
      {/* Keyed, so following a link to another search resets the fields to it. */}
      <TryForm key={`form ${key}`} languages={languages} query={query} />
      {query.q && query.language && (
        <Suspense key={`results ${key}`} fallback={<TryResultsSkeleton />}>
          <TryResults query={query} languageName={name} />
        </Suspense>
      )}
    </>
  )
}

export default function TryPage(props: PageProps<"/try">) {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex max-w-2xl flex-col gap-2">
        <h1 className="text-3xl font-semibold text-balance">Try it</h1>
        <p className="text-pretty text-muted-foreground">
          Search a dictionary through the live API. The formatted result and the JSON it came from
          are side by side.
        </p>
      </div>
      <Suspense fallback={<Skeleton className="h-32 w-full" />}>
        <Search {...props} />
      </Suspense>
    </div>
  )
}
