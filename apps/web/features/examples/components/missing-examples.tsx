import Link from "next/link"

import { ErrorCard } from "@/components/error-card"
import { NoticeCard } from "@/components/notice-card"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { apiErrorMessage } from "@/lib/api-error"
import { formatCount, percentOf } from "@/lib/format"

import { apiForVisitor } from "../../auth/session"

/**
 * A language's example coverage (LEX-5): how many verified core meanings have a verified
 * example, and those that still need one, each linking to its entry to add it.
 */
export async function MissingExamples({ language }: { language: string }) {
  const api = await apiForVisitor()
  const { data, error } = await api.GET("/v1/examples/missing", { params: { query: { language } } })
  if (!data)
    return <ErrorCard title="Coverage could not be loaded" message={apiErrorMessage(error)} />
  const { core_senses: total, with_example: done, missing } = data
  const covered = percentOf(done, total)

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col gap-2">
        <p className="text-pretty">
          <span className="font-semibold tabular-nums">
            {formatCount(done)} of {formatCount(total)}
          </span>{" "}
          verified core meanings have a verified example ({covered}%).
        </p>
        <div
          role="progressbar"
          aria-label="Meanings with an example"
          aria-valuemin={0}
          aria-valuemax={100}
          aria-valuenow={covered}
          className="h-2 overflow-hidden rounded-full bg-muted"
        >
          <div className="h-full bg-primary" style={{ width: `${covered}%` }} />
        </div>
      </div>
      {missing.length === 0 ? (
        <NoticeCard
          title={
            total === 0 ? "No verified core meanings yet" : "Every core meaning has an example"
          }
          actions={
            <Button nativeButton={false} render={<Link href="/review" />}>
              Review candidates
            </Button>
          }
        >
          {total === 0
            ? "Examples are added once a core meaning is verified. Reviewing candidates gets there sooner."
            : "Reviewers keep checking new examples as they come in."}
        </NoticeCard>
      ) : (
        <ul className="flex flex-col divide-y">
          {missing.map((m) => (
            <li
              key={m.sense_id}
              className="flex flex-wrap items-baseline justify-between gap-2 py-2"
            >
              <Link
                href={`/entries/${m.entry_id}`}
                className="min-w-0 underline-offset-4 hover:underline"
              >
                <span lang={language} className="font-medium">
                  {m.headword}
                </span>
                <span className="text-muted-foreground">: {m.gloss_en}</span>
              </Link>
              {m.pending > 0 && (
                <span className="text-sm text-muted-foreground tabular-nums">
                  {formatCount(m.pending)} waiting for review
                </span>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

export function MissingExamplesSkeleton() {
  return <Skeleton className="h-64 w-full" />
}
