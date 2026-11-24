import type { Metadata } from "next"
import { Suspense } from "react"

import { RecheckQueue, RecheckQueueSkeleton } from "@/features/rechecks/components/recheck-queue"

export const metadata: Metadata = { title: "Re-checks" }

async function Queue({ searchParams }: PageProps<"/rechecks">) {
  return <RecheckQueue params={await searchParams} />
}

export default function RechecksPage(props: PageProps<"/rechecks">) {
  return (
    <div className="mx-auto flex max-w-2xl flex-col gap-6">
      <div className="flex flex-col gap-2">
        <h1 className="text-3xl font-semibold text-balance">Re-checks</h1>
        <p className="text-pretty text-muted-foreground">
          Each quarter 5% of the verified items in each language are checked again. A wrong one goes
          back to review, and a batch with too many wrong goes back whole.
        </p>
      </div>
      <Suspense fallback={<RecheckQueueSkeleton />}>
        <Queue {...props} />
      </Suspense>
    </div>
  )
}
