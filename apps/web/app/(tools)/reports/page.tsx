import type { Metadata } from "next"
import { Suspense } from "react"

import { ReportList, ReportListSkeleton } from "@/features/reports/components/report-list"

export const metadata: Metadata = { title: "Reports" }

async function Inbox({ searchParams }: PageProps<"/reports">) {
  const { language } = await searchParams
  return <ReportList asked={typeof language === "string" ? language : undefined} />
}

export default function ReportsPage(props: PageProps<"/reports">) {
  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6">
      <div className="flex flex-col gap-2">
        <h1 className="text-3xl font-semibold text-balance">Reports</h1>
        <p className="text-pretty text-muted-foreground">
          Meanings someone says are wrong. Confirm a report with merit to send the meaning to the
          senior reviewers as a dispute.
        </p>
      </div>
      <Suspense fallback={<ReportListSkeleton />}>
        <Inbox {...props} />
      </Suspense>
    </div>
  )
}
