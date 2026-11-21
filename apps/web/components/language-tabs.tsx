import Link from "next/link"

import { Button } from "@/components/ui/button"

/** The asked-for language when it is one of codes, otherwise the first. */
export function pickLanguage(codes: string[], asked: string | undefined): string {
  return asked && codes.includes(asked) ? asked : codes[0]!
}

/**
 * One link per language to the same page, the current one marked
 * (docs/design/ui.md). Nothing when there is only one language.
 */
export function LanguageTabs({
  languages,
  current,
  path,
}: {
  languages: { code: string; name: string }[]
  current: string
  path: string
}) {
  if (languages.length < 2) return null
  return (
    <nav className="flex flex-wrap gap-2" aria-label="Language">
      {languages.map((l) => (
        <Button
          key={l.code}
          size="sm"
          variant={l.code === current ? "secondary" : "ghost"}
          aria-current={l.code === current ? "page" : undefined}
          nativeButton={false}
          render={<Link href={`${path}?${new URLSearchParams({ language: l.code })}`} />}
        >
          {l.name}
        </Button>
      ))}
    </nav>
  )
}
