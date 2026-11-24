import type { Schemas } from "@vernacular/api-client"
import Link from "next/link"

import { cn } from "cn"

import { ErrorCard } from "@/components/error-card"
import { NoticeCard } from "@/components/notice-card"
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
import { formatCount, formatDate } from "@/lib/format"

import { getLanguages, languageName } from "../../languages/get-languages"
import { getReleases } from "../get-releases"

// Examples and phrases are hidden on phones, so the table fits.
const kinds = [
  ["entries", "Words", ""],
  ["senses", "Meanings", ""],
  ["examples", "Examples", "hidden sm:table-cell"],
  ["phrases", "Phrases", "hidden sm:table-cell"],
] as const

/** requirements, Portal "Data and releases": each release's change log and counts (PRV-4). */
export async function ReleaseList() {
  const [result, names] = await Promise.all([getReleases(), getLanguages()])
  if (!result.ok) return <ErrorCard title="Releases could not be loaded" message={result.message} />
  if (result.value.length === 0) {
    return (
      <NoticeCard
        title="No releases yet."
        actions={
          <Button nativeButton={false} render={<Link href="/try" />}>
            Try the live data
          </Button>
        }
      >
        The first release is built once reviewers have verified words. Until then, the API serves
        the live data.
      </NoticeCard>
    )
  }
  return (
    <ol className="flex flex-col gap-10">
      {result.value.map((r) => (
        <Release key={r.name} release={r} name={(code) => languageName(names, code)} />
      ))}
    </ol>
  )
}

function Release({
  release: r,
  name,
}: {
  release: Schemas["Release"]
  name: (code: string) => string
}) {
  const codes = Object.keys(r.languages).sort()
  return (
    <li id={r.name} className="flex scroll-mt-6 flex-col gap-3">
      <div className="flex flex-wrap items-baseline gap-x-3">
        <h2 className="text-xl font-semibold tabular-nums">{r.name}</h2>
        <span className="text-sm text-muted-foreground">{formatDate(r.published_at)}</span>
      </div>
      <p className="text-sm text-muted-foreground tabular-nums">
        {formatCount(r.changes.added)} added · {formatCount(r.changes.changed)} changed ·{" "}
        {formatCount(r.changes.removed)} removed
      </p>
      {r.notes && <p className="whitespace-pre-line text-pretty">{r.notes}</p>}
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Language</TableHead>
            {kinds.map(([key, label, hide]) => (
              <TableHead key={key} className={cn("text-right", hide)}>
                {label}
              </TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {codes.map((code) => (
            <TableRow key={code}>
              <TableCell className="font-medium">{name(code)}</TableCell>
              {kinds.map(([key, , hide]) => (
                <TableCell key={key} className={cn("text-right tabular-nums", hide)}>
                  {formatCount(r.languages[code]![key])}
                </TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </li>
  )
}

export function ReleaseListSkeleton() {
  return (
    <div className="flex flex-col gap-3">
      <Skeleton className="h-7 w-32" />
      <Skeleton className="h-5 w-64" />
      <Skeleton className="h-32 w-full" />
    </div>
  )
}
