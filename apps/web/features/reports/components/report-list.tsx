import Link from "next/link"

import { ErrorCard } from "@/components/error-card"
import { LanguageTabs, pickLanguage } from "@/components/language-tabs"
import { NoticeCard } from "@/components/notice-card"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { apiErrorMessage } from "@/lib/api-error"
import { formatDate } from "@/lib/format-date"

import { apiForSession } from "../../auth/session"
import { getLanguages, languageName } from "../../languages/get-languages"
import { TrustBadge } from "../../trust/components/trust-badge"
import { DecideForm } from "./decide-form"

/** The report inbox for one of the reviewer's languages, oldest first (PRV-3, US-15). */
export async function ReportList({ asked }: { asked?: string }) {
  const languages = getLanguages()
  const api = await apiForSession("/reports")
  const { data: me, error: meError } = await api.GET("/v1/me")
  if (!me)
    return <ErrorCard title="Your account could not be loaded" message={apiErrorMessage(meError)} />
  const reviewing = me.approvals.filter((a) => !a.suspended)
  if (reviewing.length === 0) {
    return (
      <NoticeCard
        title="Reports are for reviewers"
        actions={
          <Button nativeButton={false} render={<Link href="/contribute" />}>
            Contribute a word
          </Button>
        }
      >
        When someone reports a wrong meaning, a reviewer of the language checks it. Admins approve
        reviewers.
      </NoticeCard>
    )
  }
  const language = pickLanguage(
    reviewing.map((a) => a.language),
    asked,
  )
  const [{ data, error }, names] = await Promise.all([
    api.GET("/v1/reports", { params: { query: { language } } }),
    languages,
  ])
  if (!data)
    return <ErrorCard title="Reports could not be loaded" message={apiErrorMessage(error)} />
  const back = `/reports?${new URLSearchParams({ language })}`

  return (
    <div className="flex flex-col gap-4">
      <LanguageTabs
        languages={reviewing.map((a) => ({
          code: a.language,
          name: languageName(names, a.language),
        }))}
        current={language}
        path="/reports"
      />
      {data.data.length === 0 ? (
        <NoticeCard
          title="No open reports."
          actions={
            <Button nativeButton={false} render={<Link href="/review" />}>
              Review words
            </Button>
          }
        >
          Nobody has reported a {languageName(names, language)} meaning as wrong.
        </NoticeCard>
      ) : (
        data.data.map((r) => (
          <Card key={r.id}>
            <CardHeader>
              <CardDescription className="flex justify-between gap-4">
                <span>
                  {languageName(names, r.language)} · {r.gloss}
                </span>
                <time dateTime={r.created_at} className="tabular-nums">
                  {formatDate(r.created_at)}
                </time>
              </CardDescription>
              <CardTitle className="flex flex-wrap items-center gap-2 text-2xl text-balance">
                <span lang={r.language}>{r.headword}</span> <TrustBadge status={r.status} />
              </CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col gap-4">
              <blockquote className="border-l-2 pl-3 text-pretty">{r.reason}</blockquote>
              <DecideForm reportId={r.id} back={back} />
            </CardContent>
          </Card>
        ))
      )}
    </div>
  )
}

export function ReportListSkeleton() {
  return <Skeleton className="h-64 w-full" />
}
