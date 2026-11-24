import Link from "next/link"

import { ErrorCard } from "@/components/error-card"
import { LanguageTabs, pickLanguage } from "@/components/language-tabs"
import { NoteDecisionForm } from "@/components/note-decision-form"
import { NoticeCard } from "@/components/notice-card"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { apiAs } from "@/lib/api"
import { apiErrorMessage } from "@/lib/api-error"
import type { Params } from "@/lib/search-params"
import { candidateDetail } from "@/lib/word-class"

import { requireSession } from "../../auth/session"
import { getLanguages, languageName } from "../../languages/get-languages"
import { recordRecheck } from "../actions"
import { readRechecks, recheckUrl } from "../recheck-url"

const review = (
  <Button nativeButton={false} render={<Link href="/review" />}>
    Review words
  </Button>
)

/**
 * The next verified item drawn for this quarter's re-check, for a senior reviewer
 * (DQ-6, US-28): still right, or wrong and back to review with a note.
 */
export async function RecheckQueue({ params }: { params: Params }) {
  const asked = readRechecks(params)
  const languages = getLanguages()
  const { session } = await requireSession("/rechecks")
  const api = apiAs(session.token)
  const next = (language: string, skip: string[]) =>
    api.GET("/v1/rechecks/next", {
      params: { query: { language, skip: skip.length > 0 ? skip : undefined } },
    })
  // When the URL names a language, fetch its item while checking the caller is
  // senior there; the API refuses languages they are not.
  const early = asked.language ? next(asked.language, asked.skip) : undefined
  early?.catch(() => undefined) // unused when the language is not theirs; awaited below otherwise
  const { data: me, error: meError } = await api.GET("/v1/me")
  if (!me)
    return <ErrorCard title="Your account could not be loaded" message={apiErrorMessage(meError)} />
  const senior = me.approvals
    .filter((a) => a.level === "senior" && !a.suspended)
    .map((a) => a.language)
  if (senior.length === 0) {
    return (
      <NoticeCard title="Re-checks are for senior reviewers" actions={review}>
        Each quarter a sample of verified words is checked again by a senior reviewer for the
        language. Admins appoint senior reviewers.
      </NoticeCard>
    )
  }
  const language = pickLanguage(senior, asked.language)
  const skip = language === asked.language ? asked.skip : []
  const [{ data, error, response }, names] = await Promise.all([
    early && language === asked.language ? early : next(language, skip),
    languages,
  ])
  const name = languageName(names, language)
  const tabs = (
    <LanguageTabs
      languages={senior.map((code) => ({ code, name: languageName(names, code) }))}
      current={language}
      path="/rechecks"
    />
  )

  if (response.status === 204) {
    return (
      <>
        {tabs}
        <NoticeCard title={`Nothing left to re-check in ${name}`} actions={review}>
          {skip.length > 0
            ? "You skipped the rest; another senior reviewer can take them."
            : "New items are drawn at the start of each quarter."}
        </NoticeCard>
      </>
    )
  }
  if (!data)
    return <ErrorCard title="Re-checks could not be loaded" message={apiErrorMessage(error)} />
  const { candidate: c } = data
  const detail = candidateDetail(c)

  return (
    <>
      {tabs}
      <Card>
        <CardHeader>
          <CardDescription>
            {name} · {data.quarter} re-check
          </CardDescription>
          <CardTitle className="text-2xl text-balance" lang={language}>
            {c.headword}
          </CardTitle>
          {detail && <p className="text-sm text-pretty text-muted-foreground">{detail}</p>}
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          <p className="text-sm text-pretty text-muted-foreground">
            This item is verified. Is it still right for “{data.gloss_en}”?
          </p>
          <NoteDecisionForm
            key={c.item_id}
            id={c.item_id}
            name="outcome"
            choices={[
              { value: "pass", label: "Still right" },
              { value: "fail", label: "Wrong: send back to review" },
            ]}
            noteHint="Needed when it is wrong: say what is wrong."
            action={recordRecheck.bind(null, c.item_id, language)}
          />
          <Button
            variant="ghost"
            size="sm"
            className="self-start"
            nativeButton={false}
            render={<Link href={recheckUrl(language, [...skip, c.item_id])} />}
          >
            Skip
          </Button>
        </CardContent>
      </Card>
    </>
  )
}

export function RecheckQueueSkeleton() {
  return <Skeleton className="h-72 w-full" />
}
