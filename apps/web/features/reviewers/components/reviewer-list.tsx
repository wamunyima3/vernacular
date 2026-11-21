import type { Schemas } from "@vernacular/api-client"
import Link from "next/link"

import { ConfirmAction } from "@/components/confirm-action"
import { ErrorCard } from "@/components/error-card"
import { NoticeCard } from "@/components/notice-card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { apiAs } from "@/lib/api"
import { apiErrorMessage } from "@/lib/api-error"

import { requireSession } from "../../auth/session"
import { getLanguages, languageName } from "../../languages/get-languages"
import { retrain } from "../actions"
import { agreementThreshold as threshold } from "../threshold"

const agreement = (r: Schemas["Reviewer"]) =>
  r.decided === 0 ? "–" : `${Math.round((100 * r.agreed) / r.decided)}%`

/** A language's reviewers with their agreement rates, for its senior reviewers (US-24). */
export async function ReviewerList({ asked }: { asked?: string }) {
  const languages = getLanguages()
  const { session } = await requireSession("/reviewers")
  const api = apiAs(session.token)
  const { data: me, error: meError } = await api.GET("/v1/me")
  if (!me)
    return <ErrorCard title="Your account could not be loaded" message={apiErrorMessage(meError)} />
  const senior = me.approvals.filter((a) => a.level === "senior" && !a.suspended)
  if (senior.length === 0) {
    return (
      <NoticeCard title="Reviewer agreement is for senior reviewers">
        Senior reviewers watch how often each reviewer agrees with settled outcomes, and retrain
        those who fall below {threshold.percent}%.
      </NoticeCard>
    )
  }
  const language = senior.some((a) => a.language === asked) ? asked! : senior[0]!.language
  const [{ data, error }, names] = await Promise.all([
    api.GET("/v1/reviewers", { params: { query: { language } } }),
    languages,
  ])
  if (!data)
    return <ErrorCard title="Reviewers could not be loaded" message={apiErrorMessage(error)} />

  return (
    <div className="flex flex-col gap-4">
      {senior.length > 1 && (
        <nav className="flex flex-wrap gap-2" aria-label="Language">
          {senior.map((a) => (
            <Button
              key={a.language}
              size="sm"
              variant={a.language === language ? "secondary" : "ghost"}
              nativeButton={false}
              render={<Link href={`/reviewers?${new URLSearchParams({ language: a.language })}`} />}
            >
              {languageName(names, a.language)}
            </Button>
          ))}
        </nav>
      )}
      <p className="text-sm text-pretty text-muted-foreground">
        {languageName(names, language)}: agreement with settled outcomes since each reviewer&apos;s
        last retraining. Below {threshold.percent}% over at least {threshold.decisions} decisions,
        approval is paused.
      </p>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead className="hidden sm:table-cell">Level</TableHead>
            <TableHead className="hidden text-right sm:table-cell">Settled</TableHead>
            <TableHead className="text-right">Agreement</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>
              <span className="sr-only">Actions</span>
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {data.reviewers.map((r) => (
            <TableRow key={r.user_id}>
              <TableCell className="max-w-40 truncate font-medium">{r.name}</TableCell>
              <TableCell className="hidden capitalize sm:table-cell">{r.level}</TableCell>
              <TableCell className="hidden text-right tabular-nums sm:table-cell">
                {r.decided}
              </TableCell>
              <TableCell className="text-right tabular-nums">{agreement(r)}</TableCell>
              <TableCell>
                <Badge variant={r.suspended ? "destructive" : "secondary"}>
                  {r.suspended ? "Paused" : "Active"}
                </Badge>
              </TableCell>
              <TableCell className="text-right">
                {r.suspended && r.user_id !== me.id && (
                  <ConfirmAction
                    label="Mark retrained"
                    title={`Mark ${r.name} retrained?`}
                    description="Their approval resumes and their agreement count starts again from zero."
                    action={retrain.bind(null, r.user_id, language)}
                    className="flex flex-col items-end gap-2"
                  />
                )}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}

export function ReviewerListSkeleton() {
  return <Skeleton className="h-48 w-full" />
}
