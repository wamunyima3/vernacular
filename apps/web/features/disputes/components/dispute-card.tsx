import type { Schemas } from "@vernacular/api-client"

import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { wordClassLabel } from "@/lib/word-class"

import { ResolveForm } from "./resolve-form"

/** DQ-8: disputes are settled within 14 days. */
const deadlineDays = 14

const decisionLabels: Record<Schemas["ReviewDecision"], string> = {
  accept: "Accepted",
  reject: "Rejected",
  correct: "Wrote another word",
}

/**
 * docs/design/ui.md, Disputes: each reviewer's decision side by side, and two
 * actions that both require a note.
 */
export function DisputeCard({
  dispute,
  languageName,
  daysOpen,
}: {
  dispute: Schemas["Dispute"]
  languageName: string
  daysOpen: number
}) {
  const { candidate } = dispute

  return (
    <Card>
      <CardHeader>
        <CardDescription className="flex justify-between gap-4">
          <span>
            {languageName} · {dispute.gloss_en}
          </span>
          <span
            className={`tabular-nums ${daysOpen > deadlineDays ? "font-medium text-destructive" : ""}`}
          >
            {daysOpen === 1 ? "1 day" : `${daysOpen} days`} open
          </span>
        </CardDescription>
        <CardTitle className="text-2xl" lang={dispute.language}>
          {candidate.headword}
          {wordClassLabel(candidate) && (
            <span className="ml-2 text-sm font-normal text-muted-foreground">
              {wordClassLabel(candidate)}
            </span>
          )}
        </CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-6">
        <ul className="grid gap-3 sm:grid-cols-2">
          {dispute.reviews.map((r, i) => (
            <li key={r.created_at + i} className="flex flex-col gap-1 rounded-lg border p-3">
              <span className="text-sm text-muted-foreground">Reviewer {i + 1}</span>
              <Badge variant={r.decision === "accept" ? "secondary" : "outline"}>
                {decisionLabels[r.decision]}
              </Badge>
              {r.correction && (
                <span lang={dispute.language} className="font-medium">
                  {r.correction}
                </span>
              )}
              {r.comment && <p className="text-sm text-pretty">{r.comment}</p>}
            </li>
          ))}
        </ul>
        <ResolveForm itemId={candidate.item_id} />
      </CardContent>
    </Card>
  )
}
