import type { Schemas } from "@vernacular/api-client"
import Link from "next/link"

import { ErrorCard } from "@/components/error-card"
import { NoticeCard } from "@/components/notice-card"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { Toaster } from "@/components/ui/sonner"
import { apiAs } from "@/lib/api"
import { apiErrorMessage } from "@/lib/api-error"

import { requireSession } from "../../auth/session"
import { getLanguages, languageName } from "../../languages/get-languages"
import { agreementThreshold } from "../../reviewers/threshold"
import { type Params, readLast, readSession, reviewUrl } from "../review-url"
import { ReviewCard } from "./review-card"
import { UndoToast } from "./undo-toast"

const aboveActionBar = { bottom: "calc(5rem + env(safe-area-inset-bottom))" }

const contribute = (
  <Button nativeButton={false} render={<Link href="/contribute" />}>
    Contribute a word
  </Button>
)

/** The current card of a review session, or why there is none, with Undo for the last decision. */
export function ReviewQueue({ params }: { params: Params }) {
  const last = readLast(params)
  return (
    <>
      {/* Above the fixed action bar, clear of the safe area. */}
      <Toaster position="bottom-center" offset={aboveActionBar} mobileOffset={aboveActionBar} />
      {last && <UndoToast last={last} params={params} />}
      <Current params={params} />
    </>
  )
}

async function Current({ params }: { params: Params }) {
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

  const approved = me.approvals.filter((a) => !a.suspended).map((a) => a.language)
  const notice = notReviewing(me.approvals)
  if (notice) return notice
  const language =
    asked.language && approved.includes(asked.language) ? asked.language : approved[0]!
  const [{ data: card, error, response }, names] = await Promise.all([
    language === asked.language && early ? early : next(language),
    languages,
  ])
  const nameOf = (code: string) => languageName(names, code)
  const current = { language, skip: asked.skip, done: asked.done }

  if (response.status === 204) {
    const others = approved
      .filter((code) => code !== language)
      .map((code) => ({ code, name: nameOf(code) }))
    return <EmptyQueue name={nameOf(language)} done={current.done} others={others} />
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

/** Why the user cannot review yet, if they cannot: no approval, or only paused ones. */
function notReviewing(approvals: Schemas["Approval"][]) {
  if (approvals.some((a) => !a.suspended)) return undefined
  if (approvals.length > 0) {
    return (
      <NoticeCard title="Your review approval is paused">
        Your agreement with settled outcomes fell below {agreementThreshold.percent}%. A senior
        reviewer for your language resumes it after retraining.
      </NoticeCard>
    )
  }
  return (
    <NoticeCard title="You are not a reviewer yet" actions={contribute}>
      Reviewers are approved per language by an admin after a short test. Contributing words is the
      way in.
    </NoticeCard>
  )
}

/** No cards left: switch to another approved language, or contribute. */
function EmptyQueue({
  name,
  done,
  others,
}: {
  name: string
  done: number
  others: { code: string; name: string }[]
}) {
  return (
    <NoticeCard
      title={`Nothing left to review in ${name}.`}
      actions={
        others.length === 0
          ? contribute
          : others.map((l) => (
              <Button
                key={l.code}
                nativeButton={false}
                render={<Link href={reviewUrl({ language: l.code, skip: [], done: 0 })} />}
              >
                Switch to {l.name}
              </Button>
            ))
      }
    >
      {done > 0
        ? `You reviewed ${done} this session. Thank you.`
        : "New words arrive from sources, contributors and AI candidates."}
    </NoticeCard>
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
