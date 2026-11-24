import type { Metadata } from "next"
import { Suspense } from "react"

import { BatchList, BatchListSkeleton } from "@/features/rechecks/components/batch-list"

export const metadata: Metadata = { title: "Re-check batches" }

async function Batches({ searchParams }: PageProps<"/rechecks/batches">) {
  const { quarter } = await searchParams
  return <BatchList quarter={typeof quarter === "string" ? quarter : undefined} />
}

export default function BatchesPage(props: PageProps<"/rechecks/batches">) {
  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6">
      <div className="flex flex-col gap-2">
        <h1 className="text-3xl font-semibold text-balance">Re-check batches</h1>
        <p className="text-pretty text-muted-foreground">
          A batch is the source and model a language&apos;s items came from. When more than 2% of a
          batch&apos;s re-checks fail, send the whole batch back to review.
        </p>
      </div>
      <Suspense fallback={<BatchListSkeleton />}>
        <Batches {...props} />
      </Suspense>
    </div>
  )
}
