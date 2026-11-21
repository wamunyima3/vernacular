import type { Metadata } from "next"
import { Suspense } from "react"

import { ReviewerList, ReviewerListSkeleton } from "@/features/reviewers/components/reviewer-list"

export const metadata: Metadata = { title: "Reviewers" }

async function List({ searchParams }: PageProps<"/reviewers">) {
  const { language } = await searchParams
  return <ReviewerList asked={typeof language === "string" ? language : undefined} />
}

export default function ReviewersPage(props: PageProps<"/reviewers">) {
  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6">
      <h1 className="text-3xl font-semibold text-balance">Reviewers</h1>
      <Suspense fallback={<ReviewerListSkeleton />}>
        <List {...props} />
      </Suspense>
    </div>
  )
}
