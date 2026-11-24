import { type ReactNode, Suspense } from "react"

import { ErrorCard } from "@/components/error-card"
import { LanguageTabs } from "@/components/language-tabs"
import { getLanguages, isActive } from "@/features/languages/get-languages"

/**
 * One tab per active language at `path`, the `language` search parameter choosing the
 * tab (the first by default), and `render` for the chosen language.
 */
export async function PerLanguage({
  searchParams,
  path,
  fallback,
  render,
}: {
  searchParams: Promise<{ language?: string | string[] }>
  path: string
  fallback: ReactNode
  render: (language: string) => ReactNode
}) {
  const [{ language: asked }, languages] = await Promise.all([searchParams, getLanguages()])
  if (!languages.ok)
    return <ErrorCard title="Languages could not be loaded" message={languages.message} />
  const active = languages.languages.filter(isActive)
  const current = active.find((l) => l.code === asked) ?? active[0]
  if (!current) return null
  return (
    <>
      <LanguageTabs languages={active} current={current.code} path={path} />
      <Suspense key={current.code} fallback={fallback}>
        {render(current.code)}
      </Suspense>
    </>
  )
}
