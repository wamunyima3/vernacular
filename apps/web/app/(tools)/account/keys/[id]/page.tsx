import type { Metadata } from "next"
import { Suspense } from "react"

import { Skeleton } from "@/components/ui/skeleton"
import { KeyUsage } from "@/features/keys/components/key-usage"

export const metadata: Metadata = { title: "Key usage" }

export default function KeyUsagePage({ params }: PageProps<"/account/keys/[id]">) {
  return (
    <div className="max-w-2xl">
      <Suspense fallback={<Skeleton className="h-96 w-full rounded-xl" />}>
        <KeyUsage params={params} />
      </Suspense>
    </div>
  )
}
