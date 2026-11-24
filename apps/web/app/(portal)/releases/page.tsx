import type { Metadata } from "next"
import { Suspense } from "react"

import { ReleaseList, ReleaseListSkeleton } from "@/features/releases/components/release-list"

export const metadata: Metadata = {
  title: "Releases",
  description: "Numbered releases of the verified data, with what changed in each.",
}

export default function ReleasesPage() {
  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6">
      <div className="flex flex-col gap-2">
        <h1 className="text-3xl font-semibold text-balance">Releases</h1>
        <p className="text-pretty text-muted-foreground">
          A release is a numbered, frozen copy of the verified words and phrases. Each lists what
          changed since the one before, and can be downloaded as open data in CSV files.
        </p>
      </div>
      <Suspense fallback={<ReleaseListSkeleton />}>
        <ReleaseList />
      </Suspense>
    </div>
  )
}
