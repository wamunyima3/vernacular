import Link from "next/link"

import { ConfirmAction } from "@/components/confirm-action"
import { ErrorCard } from "@/components/error-card"
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
import { revokeKey } from "../actions"

const lastUsed = (iso?: string) => (iso ? formatDate(iso) : "Never")

/** The signed-in user's API keys, live first, each revocable (US-01). */
export async function KeyList() {
  const api = await apiForSession("/account")
  const { data, error } = await api.GET("/v1/keys")
  if (!data) return <ErrorCard title="Keys could not be loaded" message={apiErrorMessage(error)} />
  if (data.data.length === 0) {
    return (
      <p className="text-sm text-pretty text-muted-foreground">
        No keys yet. Name your first key above to start calling the API.
      </p>
    )
  }
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Key</TableHead>
          <TableHead className="hidden sm:table-cell">Last used</TableHead>
          <TableHead className="text-right">Today</TableHead>
          <TableHead className="hidden text-right md:table-cell">30 days</TableHead>
          <TableHead>
            <span className="sr-only">Action</span>
          </TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {data.data.map((k) => (
          <TableRow key={k.id} className={k.revoked_at ? "text-muted-foreground" : undefined}>
            <TableCell className="max-w-40 truncate">
              {k.revoked_at ? (
                k.name
              ) : (
                <Link href={`/account/keys/${k.id}`} className="underline-offset-4 hover:underline">
                  {k.name}
                </Link>
              )}
            </TableCell>
            <TableCell className="font-mono">{k.prefix}…</TableCell>
            <TableCell className="hidden tabular-nums sm:table-cell">
              {lastUsed(k.last_used_at)}
            </TableCell>
            <TableCell className="text-right tabular-nums">{formatCount(k.calls_today)}</TableCell>
            <TableCell className="hidden text-right tabular-nums md:table-cell">
              {formatCount(k.calls_30d)}
            </TableCell>
            <TableCell className="text-right">
              {k.revoked_at ? (
                <Badge variant="outline">Revoked</Badge>
              ) : (
                <ConfirmAction
                  label="Revoke"
                  title={`Revoke ${k.name}?`}
                  description="Calls made with this key will fail from now on. This cannot be undone."
                  action={revokeKey.bind(null, k.id)}
                  destructive
                />
              )}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}

export function KeyListSkeleton() {
  return <Skeleton className="h-32 w-full" />
}
