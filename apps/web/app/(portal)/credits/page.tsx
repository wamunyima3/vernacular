import type { Metadata } from "next"
import { Suspense } from "react"

import { ErrorCard } from "@/components/error-card"
import { LanguageTabs } from "@/components/language-tabs"
import { CreditList, CreditListSkeleton } from "@/features/credits/components/credit-list"
import { getLanguages } from "@/features/languages/get-languages"

export const metadata: Metadata = {
  title: "Contributors",
  description: "The volunteers who contribute and review the words in Vernacular, per language.",
}

/** One tab per active language; the list is per language (requirements, Portal). */
async function Credits({ searchParams }: PageProps<"/credits">) {
  const [{ language: asked }, languages] = await Promise.all([searchParams, getLanguages()])
  if (!languages.ok)
    return <ErrorCard title="Languages could not be loaded" message={languages.message} />
  const active = languages.languages.filter((l) => l.status === "active")
  const current = active.find((l) => l.code === asked) ?? active[0]
  if (!current) return null
  return (
    <>
      <LanguageTabs languages={active} current={current.code} path="/credits" />
      <Suspense key={current.code} fallback={<CreditListSkeleton />}>
        <CreditList language={current.code} />
      </Suspense>
    </>
  )
}

export default function CreditsPage(props: PageProps<"/credits">) {
  return (
    <div className="mx-auto flex max-w-2xl flex-col gap-6">
      <div className="flex flex-col gap-2">
        <h1 className="text-3xl font-semibold text-balance">Contributors</h1>
        <p className="text-pretty text-muted-foreground">
          Volunteers contribute words and check each other&apos;s work. Everyone here helped an item
          reach the verified dataset, or kept a wrong one out.
        </p>
      </div>
      <Suspense fallback={<CreditListSkeleton />}>
        <Credits {...props} />
      </Suspense>
    </div>
  )
}
