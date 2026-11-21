import type { Metadata } from "next"
import { Suspense } from "react"

import { DisputeList, DisputeListSkeleton } from "@/features/disputes/components/dispute-list"

export const metadata: Metadata = { title: "Disputes" }

export default function DisputesPage() {
  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6">
      <h1 className="text-3xl font-semibold text-balance">Disputes</h1>
      <Suspense fallback={<DisputeListSkeleton />}>
        <DisputeList />
      </Suspense>
    </div>
  )
}
