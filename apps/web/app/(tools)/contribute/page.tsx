import type { Metadata } from "next"
import { Suspense } from "react"

import { ErrorCard } from "@/components/error-card"
import { Skeleton } from "@/components/ui/skeleton"
import { requireSession } from "@/features/auth/session"
import { ContributeForm } from "@/features/contribute/components/contribute-form"
import { correctionHref, correctingFrom } from "@/features/contribute/correcting"
import { getLanguages, isActive } from "@/features/languages/get-languages"

export const metadata: Metadata = { title: "Contribute" }

async function Contribute({
  searchParams,
}: {
  searchParams: PageProps<"/contribute">["searchParams"]
}) {
  // Fetch languages while the session is checked; a redirect still wins.
  const languages = getLanguages()
  const correcting = correctingFrom(await searchParams)
  await requireSession(correcting ? correctionHref(correcting) : "/contribute")
  const result = await languages
  if (!result.ok)
    return <ErrorCard title="Languages could not be loaded" message={result.message} />
  return (
    <>
      {correcting && (
        <p className="text-pretty">
          Correcting{" "}
          <span lang={correcting.language} className="font-medium">
            {correcting.headword}
          </span>
          : {correcting.gloss_en}. Change what is wrong; the current meaning stays until reviewers
          verify yours.
        </p>
      )}
      <ContributeForm
        key={correcting?.corrects}
        languages={result.languages.filter(isActive)}
        correcting={correcting}
      />
    </>
  )
}

export default function ContributePage({ searchParams }: PageProps<"/contribute">) {
  return (
    <div className="flex max-w-xl flex-col gap-6">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-semibold text-balance">Contribute a word</h1>
        <p className="text-pretty text-muted-foreground">
          Words you add are checked by two native-speaker reviewers before anyone relies on them.
        </p>
      </div>
      <Suspense fallback={<Skeleton className="h-96 w-full rounded-xl" />}>
        <Contribute searchParams={searchParams} />
      </Suspense>
    </div>
  )
}
