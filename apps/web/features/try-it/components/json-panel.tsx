import { CopyButton } from "@/components/copy-button"

/** The raw response beside the formatted result (US-02, docs/design/ui.md). */
export function JsonPanel({ json, request }: { json: string; request: string }) {
  return (
    <div className="flex min-w-0 flex-col gap-2 rounded-lg border bg-muted/40 p-3">
      <div className="flex items-center justify-between gap-2">
        <code className="truncate text-xs text-muted-foreground">{request}</code>
        <CopyButton text={json} label="Copy JSON" size="sm" />
      </div>
      <pre className="max-h-[32rem] overflow-auto text-xs leading-relaxed">
        <code>{json}</code>
      </pre>
    </div>
  )
}
