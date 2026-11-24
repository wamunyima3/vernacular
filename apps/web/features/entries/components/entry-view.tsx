import type { Schemas } from "@vernacular/api-client"
import type { ReactNode } from "react"

import Link from "next/link"

import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { wordClassLabel } from "@/lib/word-class"

import { AddExample } from "../../contribute/components/add-example"
import { correctionHref } from "../../contribute/correcting"
import { TrustBadge } from "../../trust/components/trust-badge"
import { ReportDialog } from "./report-dialog"

const sourceKey = (src: Schemas["SourceRef"]) => `${src.name}|${src.model ?? ""}`

/** Every source cited by the entry or its meanings and examples, once each. */
function sourcesOf(entry: Schemas["Entry"]): Schemas["SourceRef"][] {
  const all = [
    entry.trust,
    ...entry.senses.flatMap((s) => [s.trust, ...s.examples.map((x) => x.trust)]),
  ]
  const seen = new Map<string, Schemas["SourceRef"]>()
  for (const src of all.flatMap((t) => t.sources)) seen.set(sourceKey(src), src)
  return [...seen.values()]
}

/**
 * A full entry (US-16, docs/design/ui.md "Entry"): the word's facts on one line, numbered
 * meanings with examples and trust, then variants and sources.
 */
export function EntryView({
  entry,
  languageName,
  pronunciation,
}: {
  entry: Schemas["Entry"]
  languageName: string
  /** The entry's recordings, streamed in after the entry (features/audio). */
  pronunciation?: ReactNode
}) {
  const facts = [
    wordClassLabel(entry),
    entry.plural && `plural ${entry.plural}`,
    entry.root && `root ${entry.root}`,
    entry.variety && `variety ${entry.variety}`,
    entry.register !== "neutral" && entry.register,
  ].filter(Boolean)
  const sources = sourcesOf(entry)

  return (
    <article className="flex max-w-3xl flex-col gap-8">
      <header className="flex flex-col gap-2">
        <p className="text-sm text-muted-foreground">{languageName}</p>
        <h1 className="text-4xl font-semibold text-balance" lang={entry.language}>
          {entry.headword}
        </h1>
        {facts.length > 0 && (
          <p className="text-pretty text-muted-foreground">{facts.join(" · ")}</p>
        )}
      </header>
      {pronunciation}

      <ol className="flex flex-col gap-6">
        {entry.senses.map((s, i) => (
          <li key={s.id} className="flex gap-3">
            <span className="text-muted-foreground tabular-nums">{i + 1}.</span>
            <div className="flex min-w-0 flex-1 flex-col gap-2">
              <div className="flex flex-wrap items-center gap-2">
                <span className="text-xl font-medium">{s.gloss}</span>
                <TrustBadge status={s.trust.status} />
              </div>
              {s.definition && <p className="text-pretty">{s.definition}</p>}
              {s.examples.length > 0 && (
                <ul className="flex flex-col gap-2 border-l-2 pl-3">
                  {s.examples.map((x) => (
                    <li key={x.id} className="flex flex-col">
                      <span lang={entry.language} className="text-pretty italic">
                        {x.text}
                      </span>
                      <span className="text-sm text-pretty text-muted-foreground">
                        {x.translation_en}
                      </span>
                    </li>
                  ))}
                </ul>
              )}
              <div className="flex flex-wrap gap-1">
                <ReportDialog senseId={s.id} meaning={`${entry.headword}: ${s.gloss}`} />
                <Button
                  variant="ghost"
                  size="sm"
                  nativeButton={false}
                  render={
                    <Link
                      href={correctionHref({
                        corrects: s.id,
                        language: entry.language,
                        headword: entry.headword,
                        gloss_en: s.gloss,
                      })}
                    />
                  }
                >
                  Suggest a correction
                </Button>
              </div>
              <AddExample entryId={entry.id} senseId={s.id} />
            </div>
          </li>
        ))}
      </ol>

      <Separator />

      <dl className="grid gap-4 text-sm sm:grid-cols-[8rem_1fr]">
        {entry.variants.length > 0 && (
          <>
            <dt className="text-muted-foreground">Variants</dt>
            <dd className="flex flex-col gap-1">
              {entry.variants.map((v) => (
                <span key={v.spelling}>
                  <span lang={entry.language}>{v.spelling}</span>
                  {v.note && <span className="text-muted-foreground"> ({v.note})</span>}
                </span>
              ))}
            </dd>
          </>
        )}
        <dt className="text-muted-foreground">Sources</dt>
        <dd className="flex flex-col gap-1">
          {sources.map((src) => (
            <span key={sourceKey(src)}>
              {src.url ? (
                <a href={src.url} className="underline underline-offset-4 hover:text-foreground">
                  {src.name}
                </a>
              ) : (
                src.name
              )}
              <span className="text-muted-foreground"> · {src.model ?? src.licence}</span>
            </span>
          ))}
        </dd>
      </dl>
    </article>
  )
}
