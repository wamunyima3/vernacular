import type { Metadata } from "next"
import { Suspense } from "react"

import { ErrorCard } from "@/components/error-card"
import { Skeleton } from "@/components/ui/skeleton"
import { requireSession } from "@/features/auth/session"
import { getLanguages, isActive } from "@/features/languages/get-languages"
import { QualifyFlow } from "@/features/qualify/components/qualify-flow"

export const metadata: Metadata = { title: "Reviewer test" }

async function Qualify() {
  const languages = getLanguages()
  await requireSession("/qualify")
  const result = await languages
  if (!result.ok)
    return <ErrorCard title="Languages could not be loaded" message={result.message} />
  return <QualifyFlow languages={result.languages.filter(isActive)} />
}

export default function QualifyPage() {
  return (
    <div className="flex max-w-xl flex-col gap-6">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-semibold text-balance">Become a reviewer</h1>
        <p className="text-pretty text-muted-foreground">
          Reviewers are native speakers. Ten questions from verified words check that you know the
          language; a senior reviewer then confirms you.
        </p>
      </div>
      <Suspense fallback={<Skeleton className="h-40 w-full rounded-xl" />}>
        <Qualify />
      </Suspense>
    </div>
  )
}
