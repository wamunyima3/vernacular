import type { Metadata } from "next"

import { errorCodes } from "@/features/docs/error-codes"

export const metadata: Metadata = { title: "Errors" }

/** Target of every error body's docs_url: /docs/errors#<code>. */
export default function ErrorsPage() {
  return (
    <article className="flex max-w-3xl flex-col gap-8">
      <header className="flex flex-col gap-3">
        <h1 className="text-3xl font-semibold text-balance">Errors</h1>
        <p className="text-pretty text-muted-foreground">
          Every error response has the same body. Match on <code className="font-mono">code</code>,
          which never changes; <code className="font-mono">message</code> is for people and may
          change. New codes can be added, so treat an unknown code as a generic error.
        </p>
        <pre className="overflow-x-auto rounded-lg bg-muted p-4 font-mono text-sm">
          {`{
  "code": "not_found",
  "message": "Route not found",
  "docs_url": "https://…/docs/errors#not_found"
}`}
        </pre>
      </header>

      <dl className="flex flex-col divide-y border-y">
        {Object.entries(errorCodes).map(([code, { status, meaning, action }]) => (
          <div key={code} id={code} className="flex scroll-mt-6 flex-col gap-1 py-4">
            <dt className="flex items-baseline gap-3">
              <code className="font-mono font-medium">{code}</code>
              <span className="text-sm text-muted-foreground tabular-nums">HTTP {status}</span>
            </dt>
            <dd className="text-pretty">{meaning}</dd>
            <dd className="text-sm text-pretty text-muted-foreground">{action}</dd>
          </div>
        ))}
      </dl>
    </article>
  )
}
