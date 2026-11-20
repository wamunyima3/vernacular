import Link from "next/link"

import { ErrorCard } from "@/components/error-card"
import { NoticeCard } from "@/components/notice-card"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { apiAs } from "@/lib/api"
import { apiErrorMessage } from "@/lib/api-error"

import { requireSession } from "../../auth/session"
import { getLanguages } from "../../languages/get-languages"
import { readSession, reviewUrl } from "../review-url"
import { ReviewCard } from "./review-card"

const contribute = (
  <Button nativeButton={false} render={<Link href="/contribute" />}>
    Contribute a word
  </Button>
)

/** The current card of a review session, or why there is none. */
export async function ReviewQueue({
  params,
}: {
  params: Record<string, string | string[] | undefined>
}) {
  const asked = readSession(params)
  const languages = getLanguages()
  const { session } = await requireSession("/review")
  const api = apiAs(session.token)
  const next = (language: string) =>
    api.GET("/v1/review/next", {
      params: { query: { language, skip: asked.skip.length > 0 ? asked.skip : undefined } },
    })
  // After the first card the URL names the language, so fetch the card while
  // checking approvals; the API refuses languages the reviewer is not approved for.
  const early = asked.language ? next(asked.language) : undefined
  early?.catch(() => undefined) // unused when the language is not approved; awaited below otherwise
  const { data: me, error: meError } = await api.GET("/v1/me")
  if (!me)
    return <ErrorCard title="Your account could not be loaded" message={apiErrorMessage(meError)} />

  const approved = me.approvals.map((a) => a.language)
  if (approved.length === 0) {
    return (
      <NoticeCard title="You are not a reviewer yet" actions={contribute}>
        Reviewers are approved per language by an admin after a short test. Contributing words is
        the way in.
      </NoticeCard>
    )
  }
  const language =
    asked.language && approved.includes(asked.language) ? asked.language : approved[0]!
  const [{ data: card, error, response }, names] = await Promise.all([
    language === asked.language && early ? early : next(language),
    languages,
  ])
  const nameOf = (code: string) =>
    (names.ok ? names.languages.find((l) => l.code === code)?.name : undefined) ?? code
  const current = { ...asked, language }

  if (response.status === 204) {
    const others = approved.filter((code) => code !== language)
    return (
      <NoticeCard
        title={`Nothing left to review in ${nameOf(language)}.`}
        actions={
          others.length === 0
            ? contribute
            : others.map((code) => (
                <Button
                  key={code}
                  nativeButton={false}
                  render={<Link href={reviewUrl({ language: code, skip: [], done: 0 })} />}
                >
                  Switch to {nameOf(code)}
                </Button>
              ))
        }
      >
        {current.done > 0
          ? `You reviewed ${current.done} this session. Thank you.`
          : "New words arrive from sources, contributors and AI candidates."}
      </NoticeCard>
    )
  }
  if (!card)
    return (
      <ErrorCard title="The review queue could not be loaded" message={apiErrorMessage(error)} />
    )
  return (
    <ReviewCard
      key={card.concept_id}
      card={card}
      languageName={nameOf(language)}
      session={current}
    />
  )
}

/** The card's shape while the next one loads. */
export function ReviewQueueSkeleton() {
  return (
    <div className="flex flex-col gap-6">
      <Skeleton className="h-4 w-40" />
      <Skeleton className="h-9 w-48" />
      <div className="flex flex-col gap-2">
        <Skeleton className="h-12 w-full" />
        <Skeleton className="h-12 w-full" />
      </div>
    </div>
  )
}
