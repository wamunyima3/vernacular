import Link from "next/link"

import { ErrorCard } from "@/components/error-card"
import { NoticeCard } from "@/components/notice-card"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { apiAs } from "@/lib/api"
import { apiErrorMessage } from "@/lib/api-error"

import { requireSession } from "../../auth/session"
import { getLanguages, languageName } from "../../languages/get-languages"
import { DisputeCard } from "./dispute-card"

const day = 24 * 60 * 60 * 1000

/** Whole days since a time; the page renders per request, so now is the request's. */
const daysSince = (iso: string) => Math.floor((Date.now() - Date.parse(iso)) / day)

const review = (
  <Button nativeButton={false} render={<Link href="/review" />}>
    Review words
  </Button>
)

/** Disputes the signed-in senior reviewer can settle, oldest first (US-23). */
export async function DisputeList() {
  const languages = getLanguages()
  const { session } = await requireSession("/disputes")
  const { data, error, response } = await apiAs(session.token).GET("/v1/disputes")
  if (response.status === 403) {
    return (
      <NoticeCard title="Disputes are for senior reviewers" actions={review}>
        When reviewers disagree, a senior reviewer for the language settles it. Admins appoint
        senior reviewers.
      </NoticeCard>
    )
  }
  if (!data)
    return <ErrorCard title="Disputes could not be loaded" message={apiErrorMessage(error)} />
  if (data.disputes.length === 0) {
    return (
      <NoticeCard title="No disputes." actions={review}>
        Reviewers agree on everything for now.
      </NoticeCard>
    )
  }
  const names = await languages
  return data.disputes.map((d) => (
    <DisputeCard
      key={d.candidate.item_id}
      dispute={d}
      languageName={languageName(names, d.language)}
      daysOpen={daysSince(d.disputed_since)}
    />
  ))
}

export function DisputeListSkeleton() {
  return (
    <div className="flex flex-col gap-4">
      <Skeleton className="h-64 w-full" />
      <Skeleton className="h-64 w-full" />
    </div>
  )
}
