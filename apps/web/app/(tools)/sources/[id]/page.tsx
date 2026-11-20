import type { Metadata } from "next"
import { Suspense } from "react"

import { Skeleton } from "@/components/ui/skeleton"
import { SourceEditor } from "@/features/sources/components/source-editor"

export const metadata: Metadata = { title: "Edit source" }

export default function SourcePage({ params }: PageProps<"/sources/[id]">) {
  return (
    <div className="max-w-xl">
      <Suspense fallback={<Skeleton className="h-96 w-full rounded-xl" />}>
        <SourceEditor params={params} />
      </Suspense>
    </div>
  )
}
