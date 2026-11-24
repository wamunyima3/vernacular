import type { Schemas } from "@vernacular/api-client"

import { ConfirmAction } from "@/components/confirm-action"
import { formatDate } from "@/lib/format"

import { confirmTest } from "../actions"

/** Passed reviewer tests in one language, for its senior reviewers to confirm (DQ-3). */
export function PendingTests({
  tests,
  languageName,
}: {
  tests: Schemas["PendingQualification"][]
  languageName: string
}) {
  if (tests.length === 0) return null
  return (
    <section className="flex flex-col gap-3">
      <h2 className="text-xl font-semibold text-balance">Passed tests to confirm</h2>
      <ul className="flex flex-col divide-y rounded-lg border">
        {tests.map((t) => (
          <li key={t.id} className="flex flex-wrap items-center justify-between gap-3 p-3">
            <div className="flex min-w-0 flex-col">
              <span className="truncate font-medium">{t.name}</span>
              <span className="text-sm text-muted-foreground tabular-nums">
                {t.score} of {t.total} · {formatDate(t.answered_at)}
              </span>
            </div>
            <ConfirmAction
              label="Confirm"
              title={`Confirm ${t.name} as a ${languageName} reviewer?`}
              description="Confirm only someone you know speaks the language natively. They can review straight away."
              action={confirmTest.bind(null, t.id, t.language)}
              className="flex flex-col items-end gap-2"
            />
          </li>
        ))}
      </ul>
    </section>
  )
}
