import type { Metadata } from "next"
import { Suspense } from "react"

import { PerLanguage } from "@/components/per-language"
import { CreditList, CreditListSkeleton } from "@/features/credits/components/credit-list"

export const metadata: Metadata = {
  title: "Contributors",
  description: "The volunteers who contribute and review the words in Vernacular, per language.",
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
        <PerLanguage
          searchParams={props.searchParams}
          path="/credits"
          fallback={<CreditListSkeleton />}
          render={(language) => <CreditList language={language} />}
        />
      </Suspense>
    </div>
  )
}
