import type { Metadata } from "next"
import Link from "next/link"
import { Suspense } from "react"

import { Button } from "@/components/ui/button"
import { DisputeList, DisputeListSkeleton } from "@/features/disputes/components/dispute-list"

export const metadata: Metadata = { title: "Disputes" }

export default function DisputesPage() {
  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6">
      <div className="flex items-center justify-between gap-4">
        <h1 className="text-3xl font-semibold text-balance">Disputes</h1>
        <Button
          variant="outline"
          size="sm"
          nativeButton={false}
          render={<Link href="/reviewers" />}
        >
          Reviewers
        </Button>
      </div>
      <Suspense fallback={<DisputeListSkeleton />}>
        <DisputeList />
      </Suspense>
    </div>
  )
}
