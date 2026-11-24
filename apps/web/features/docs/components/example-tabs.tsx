import { readFileSync } from "node:fs"
import { join } from "node:path"

import { CopyButton } from "@/components/copy-button"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

// The files in examples/, which `make examples` runs against the API, so what the docs
// show is what was run (US-02). Next runs from apps/web.
// Read once when the module loads, so the docs page stays static.
const examples = [
  { file: "search.sh", label: "curl", run: "sh search.sh" },
  { file: "search.mjs", label: "JavaScript", run: "node search.mjs" },
  { file: "search.py", label: "Python", run: "python3 search.py" },
  { file: "search.go", label: "Go", run: "go run search.go" },
].map((e) => ({ ...e, code: readFileSync(join(process.cwd(), "../../examples", e.file), "utf8") }))

/** The same request in each language, as tabs with a copy button. */
export function ExampleTabs() {
  return (
    <Tabs defaultValue={examples[0]!.file} className="my-2">
      <TabsList>
        {examples.map((e) => (
          <TabsTrigger key={e.file} value={e.file}>
            {e.label}
          </TabsTrigger>
        ))}
      </TabsList>
      {examples.map((e) => (
        <TabsContent key={e.file} value={e.file} className="flex flex-col gap-2">
          <div className="flex items-center justify-between gap-2 text-sm text-muted-foreground">
            <span>
              {e.file} · run with <code className="font-mono">{e.run}</code>
            </span>
            <CopyButton text={e.code} label={`Copy the ${e.label} example`} size="sm" />
          </div>
          <pre className="overflow-x-auto rounded-lg bg-muted p-4 font-mono text-sm [tab-size:4]">
            <code>{e.code}</code>
          </pre>
        </TabsContent>
      ))}
    </Tabs>
  )
}
