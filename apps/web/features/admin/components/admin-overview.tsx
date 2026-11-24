import { AdminOnly } from "@/components/admin-only"
import { ConfirmAction } from "@/components/confirm-action"
import { ErrorCard } from "@/components/error-card"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import { percentOf } from "@/lib/format"
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

import { sessionWithAdmin } from "../../auth/admin"
import { getLanguages, languageName } from "../../languages/get-languages"
import { BuildReleaseForm } from "../../releases/components/build-release-form"
import { revoke } from "../actions"
import { ApproveForm } from "./approve-form"

const percent = (n: number, of: number) => (of === 0 ? "–" : `${percentOf(n, of)}%`)

/** docs/design/ui.md, Admin: reviewer approvals and provider quality, both as tables (US-26). */
export async function AdminOverview() {
  const { session, isAdmin } = await sessionWithAdmin("/admin")
  if (!isAdmin) {
    return <AdminOnly>Admins approve reviewers and watch AI provider quality.</AdminOnly>
  }
  const api = apiAs(session.session.token)
  const [approvals, quality, names] = await Promise.all([
    api.GET("/v1/approvals"),
    api.GET("/v1/provider-quality"),
    getLanguages(),
  ])
  if (!approvals.data || !quality.data) {
    return (
      <ErrorCard
        title="The admin overview could not be loaded"
        message={apiErrorMessage(approvals.error ?? quality.error)}
      />
    )
  }
  const name = (code: string) => languageName(names, code)

  return (
    <div className="flex flex-col gap-10">
      <section className="flex flex-col gap-4" aria-labelledby="approvals">
        <h2 id="approvals" className="text-xl font-semibold text-balance">
          Reviewer approvals
        </h2>
        <ApproveForm languages={names.ok ? names.languages : []} />
        {approvals.data.approvals.length === 0 ? (
          <p className="text-sm text-muted-foreground">No reviewers approved yet.</p>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Reviewer</TableHead>
                <TableHead>Language</TableHead>
                <TableHead className="hidden sm:table-cell">Level</TableHead>
                <TableHead>
                  <span className="sr-only">Actions</span>
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {approvals.data.approvals.map((a) => (
                <TableRow key={`${a.user_id}-${a.language}`}>
                  <TableCell className="max-w-48">
                    <div className="truncate font-medium">{a.name}</div>
                    <div className="truncate text-muted-foreground">{a.email}</div>
                  </TableCell>
                  <TableCell>
                    {name(a.language)}
                    {/* On phones the level sits under the language, so the row fits. */}
                    <div className="text-muted-foreground capitalize sm:hidden">
                      {a.level}
                      {a.suspended && " · paused"}
                    </div>
                  </TableCell>
                  <TableCell className="hidden sm:table-cell">
                    <span className="capitalize">{a.level}</span>
                    {a.suspended && (
                      <Badge variant="destructive" className="ml-2">
                        Paused
                      </Badge>
                    )}
                  </TableCell>
                  <TableCell className="text-right">
                    <ConfirmAction
                      label="Remove"
                      title={`Remove ${a.name} as a ${name(a.language)} reviewer?`}
                      description="They can no longer review this language. Their past reviews stay."
                      action={revoke.bind(null, a.user_id, a.language)}
                      className="flex flex-col items-end gap-2"
                    />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </section>

      <section className="flex flex-col gap-4" aria-labelledby="releases">
        <h2 id="releases" className="text-xl font-semibold text-balance">
          Releases
        </h2>
        <BuildReleaseForm />
      </section>

      <section className="flex flex-col gap-4" aria-labelledby="providers">
        <h2 id="providers" className="text-xl font-semibold text-balance">
          Provider quality
        </h2>
        {quality.data.providers.length === 0 ? (
          <p className="text-sm text-muted-foreground">No AI candidates have been reviewed yet.</p>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Provider</TableHead>
                <TableHead>Language</TableHead>
                <TableHead className="text-right">Decisions</TableHead>
                <TableHead className="text-right">Accepted</TableHead>
                <TableHead className="text-right">Corrected</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {quality.data.providers.map((p) => (
                <TableRow key={`${p.provider}-${p.language}`}>
                  <TableCell className="font-medium">{p.provider}</TableCell>
                  <TableCell>{name(p.language)}</TableCell>
                  <TableCell className="text-right tabular-nums">{p.decided}</TableCell>
                  <TableCell className="text-right tabular-nums">
                    {percent(p.accepted, p.decided)}
                  </TableCell>
                  <TableCell className="text-right tabular-nums">
                    {percent(p.corrected, p.decided)}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </section>
    </div>
  )
}

export function AdminOverviewSkeleton() {
  return (
    <div className="flex flex-col gap-4">
      <Skeleton className="h-10 w-full" />
      <Skeleton className="h-48 w-full" />
    </div>
  )
}
