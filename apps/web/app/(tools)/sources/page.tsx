import type { Metadata } from "next"
import { Suspense } from "react"

import { Skeleton } from "@/components/ui/skeleton"
import { SourcesOverview } from "@/features/sources/components/sources-overview"

export const metadata: Metadata = { title: "Sources" }

export default function SourcesPage() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex max-w-2xl flex-col gap-2">
        <h1 className="text-2xl font-semibold text-balance">Source register</h1>
        <p className="text-pretty text-muted-foreground">
          Every source the dataset draws on, with its licence. Nothing is imported from a source
          that is not listed here.
        </p>
      </div>
      <Suspense fallback={<Skeleton className="h-64 w-full rounded-xl" />}>
        <SourcesOverview />
      </Suspense>
    </div>
  )
}
