import type { Metadata } from "next"
import { Suspense } from "react"

import { ErrorCard } from "@/components/error-card"
import { Skeleton } from "@/components/ui/skeleton"
import { requireSession } from "@/features/auth/session"
import { ContributeForm } from "@/features/contribute/components/contribute-form"
import { getLanguages } from "@/features/languages/get-languages"

export const metadata: Metadata = { title: "Contribute" }

async function Contribute() {
  // Fetch languages while the session is checked; a redirect still wins.
  const languages = getLanguages()
  await requireSession("/contribute")
  const result = await languages
  if (!result.ok)
    return <ErrorCard title="Languages could not be loaded" message={result.message} />
  return <ContributeForm languages={result.languages.filter((l) => l.status === "active")} />
}

export default function ContributePage() {
  return (
    <div className="flex max-w-xl flex-col gap-6">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-semibold text-balance">Contribute a word</h1>
        <p className="text-pretty text-muted-foreground">
          Words you add are checked by two native-speaker reviewers before anyone relies on them.
        </p>
      </div>
      <Suspense fallback={<Skeleton className="h-96 w-full rounded-xl" />}>
        <Contribute />
      </Suspense>
    </div>
  )
}
