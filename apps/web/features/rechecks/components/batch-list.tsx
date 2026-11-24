import type { Schemas } from "@vernacular/api-client"

import { AdminOnly } from "@/components/admin-only"
import { ConfirmAction } from "@/components/confirm-action"
import { ErrorCard } from "@/components/error-card"
import { NoticeCard } from "@/components/notice-card"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
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
import { drawSample, sendBackBatch } from "../actions"

/**
 * A quarter's re-checks per batch, worst first (DQ-6, US-28): a batch failing
 * more than 2% of its checks can be sent back to review, once.
 */
export async function BatchList({ quarter }: { quarter?: string }) {
  const api = await apiForSession("/rechecks/batches")
  const { data, error, response } = await api.GET("/v1/rechecks/batches", {
    params: { query: { quarter } },
  })
  if (response.status === 403)
    return (
      <AdminOnly>
        Admins see each quarter&apos;s re-check results and send back failing batches.
      </AdminOnly>
    )
  if (!data)
    return (
      <ErrorCard title="Re-check batches could not be loaded" message={apiErrorMessage(error)} />
    )

  if (data.batches.length === 0) {
    return (
      <NoticeCard
        title={`Nothing drawn for ${data.quarter} yet`}
        actions={
          <ConfirmAction
            label="Draw the sample now"
            title={`Draw the ${data.quarter} sample?`}
            description="5% of each language's verified items go to senior reviewers to check again. The worker does this at the start of each quarter; a quarter is drawn once."
            action={drawSample}
          />
        }
      >
        The worker draws each quarter&apos;s sample when the quarter starts.
      </NoticeCard>
    )
  }
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Batch</TableHead>
          <TableHead className="text-right">Drawn</TableHead>
          <TableHead className="text-right">Checked</TableHead>
          <TableHead className="text-right">Failed</TableHead>
          <TableHead>
            <span className="sr-only">Action</span>
          </TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {data.batches.map((b) => (
          <TableRow key={`${b.language}/${b.source_id}/${b.model}`}>
            <TableCell className="whitespace-normal">
              <span className="font-medium">{b.source}</span>
              {b.model && <span className="text-muted-foreground"> · {b.model}</span>}
              <span className="block text-sm text-muted-foreground">{b.language}</span>
            </TableCell>
            <TableCell className="text-right tabular-nums">{formatCount(b.sampled)}</TableCell>
            <TableCell className="text-right tabular-nums">{formatCount(b.checked)}</TableCell>
            <TableCell className="text-right tabular-nums">
              {b.fails ? (
                <Badge variant="destructive">{formatCount(b.failed)}</Badge>
              ) : (
                formatCount(b.failed)
              )}
            </TableCell>
            <TableCell className="text-right">
              <BatchAction quarter={data.quarter} batch={b} />
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}

function BatchAction({ quarter, batch: b }: { quarter: string; batch: Schemas["RecheckBatch"] }) {
  if (b.sent_back_at) {
    return (
      <span className="text-sm text-muted-foreground tabular-nums">
        {formatCount(b.sent_back_items ?? 0)} sent back {formatDate(b.sent_back_at)}
      </span>
    )
  }
  if (!b.fails) return null
  const name = b.model ? `${b.source} · ${b.model}` : b.source
  return (
    <ConfirmAction
      label="Send back"
      title={`Send ${name} back to review?`}
      description={`More than 2% of its checks failed. Every verified ${b.language} item from this batch goes back to the review queue, where only new reviews count. This cannot be undone.`}
      action={sendBackBatch.bind(null, {
        quarter,
        language: b.language,
        source_id: b.source_id,
        model: b.model,
      })}
      destructive
    />
  )
}

export function BatchListSkeleton() {
  return <Skeleton className="h-64 w-full" />
}
