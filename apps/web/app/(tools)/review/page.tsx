import type { Metadata } from "next"
import { Suspense } from "react"

import { ReviewQueue, ReviewQueueSkeleton } from "@/features/review/components/review-queue"

export const metadata: Metadata = { title: "Review" }

async function Queue({ searchParams }: PageProps<"/review">) {
  return <ReviewQueue params={await searchParams} />
}

export default function ReviewPage(props: PageProps<"/review">) {
  return (
    <div className="mx-auto flex max-w-xl flex-col gap-6">
      <Suspense fallback={<ReviewQueueSkeleton />}>
        <Queue {...props} />
      </Suspense>
    </div>
  )
}
