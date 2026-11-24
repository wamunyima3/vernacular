import Link from "next/link"
import { notFound } from "next/navigation"

import { ErrorCard } from "@/components/error-card"
import { NoticeCard } from "@/components/notice-card"
import { Button } from "@/components/ui/button"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { apiErrorMessage } from "@/lib/api-error"
import { formatCount, formatDate } from "@/lib/format"

import { apiForSession } from "../../auth/session"
import { usageSummary } from "../usage"

function Stat({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex flex-col gap-1 rounded-lg border p-3">
      <dt className="text-sm text-muted-foreground">{label}</dt>
      <dd className="text-xl font-semibold tabular-nums">{value}</dd>
    </div>
  )
}

/** One key's calls over the last 30 days, by day, with its rate limit (DEV-5, US-13). */
export async function KeyUsage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params
  const api = await apiForSession(`/account/keys/${id}`)
  const usage = await api.GET("/v1/keys/{id}/usage", { params: { path: { id } } })
  if (usage.response.status === 404 || usage.response.status === 400) notFound()
  if (!usage.data) {
    return <ErrorCard title="Usage could not be loaded" message={apiErrorMessage(usage.error)} />
  }
  const { name, prefix, days, limit_per_minute: limit } = usage.data
  const s = usageSummary(days)

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1">
        <Link
          href="/account"
          className="text-sm text-muted-foreground underline-offset-4 hover:underline"
        >
          ← Account
        </Link>
        <h1 className="text-2xl font-semibold text-balance">
          {name} <span className="font-mono text-base">{prefix}…</span>
        </h1>
        <p className="text-sm text-pretty text-muted-foreground">
          Calls made with this key, by UTC day. New calls show here within a minute.
        </p>
      </div>
      <dl className="grid grid-cols-2 gap-3 sm:grid-cols-4">
        <Stat label="Today" value={formatCount(s.today)} />
        <Stat label="Last 30 days" value={formatCount(s.calls)} />
        <Stat label="Errors" value={formatCount(s.errors)} />
        <Stat label="Over the limit" value={formatCount(s.limited)} />
      </dl>
      <p className="text-sm text-pretty text-muted-foreground">
        {limit > 0
          ? `Each key may make ${formatCount(limit)} calls a minute. Calls over that get 429; wait until the time in X-RateLimit-Reset.`
          : "Keys have no rate limit on this server."}
      </p>
      {days.length === 0 ? (
        <NoticeCard
          title="No calls in the last 30 days"
          actions={
            <Button nativeButton={false} render={<Link href="/docs" />}>
              Read the quickstart
            </Button>
          }
        >
          Send the key in the X-API-Key header and its calls show here.
        </NoticeCard>
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Day</TableHead>
              <TableHead className="w-1/2">
                <span className="sr-only">Share of the busiest day</span>
              </TableHead>
              <TableHead className="text-right">Calls</TableHead>
              <TableHead className="text-right">Errors</TableHead>
              <TableHead className="text-right">429</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {days.map((d) => (
              <TableRow key={d.day}>
                <TableCell className="tabular-nums">{formatDate(`${d.day}T12:00:00Z`)}</TableCell>
                <TableCell aria-hidden>
                  <div
                    className="h-2 rounded-full bg-primary"
                    style={{ width: `${Math.max(2, (d.calls / s.busiest) * 100)}%` }}
                  />
                </TableCell>
                <TableCell className="text-right tabular-nums">{formatCount(d.calls)}</TableCell>
                <TableCell className="text-right tabular-nums">{formatCount(d.errors)}</TableCell>
                <TableCell className="text-right tabular-nums">{formatCount(d.limited)}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}
    </div>
  )
}
