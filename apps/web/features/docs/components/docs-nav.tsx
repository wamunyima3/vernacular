"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"

import { cn } from "cn"

export type DocsLink = { href: string; title: string }

function Links({ links }: { links: DocsLink[] }) {
  const pathname = usePathname()
  return (
    <ul className="flex flex-col gap-1 text-sm">
      {links.map((l) => {
        const current = pathname === l.href
        return (
          <li key={l.href}>
            <Link
              href={l.href}
              aria-current={current ? "page" : undefined}
              className={cn(
                "block rounded-md px-2 py-1.5 hover:bg-muted",
                current ? "bg-muted font-medium" : "text-muted-foreground",
              )}
            >
              {l.title}
            </Link>
          </li>
        )
      })}
    </ul>
  )
}

/** docs/design/ui.md "Docs": a left nav on wide screens, a collapsible list on phones. */
export function DocsNav({ links }: { links: DocsLink[] }) {
  return (
    <>
      <nav aria-label="Docs" className="hidden lg:block">
        <Links links={links} />
      </nav>
      <details className="rounded-lg border lg:hidden">
        <summary className="cursor-pointer px-3 py-2 text-sm font-medium">Docs pages</summary>
        <nav aria-label="Docs" className="px-1 pb-2">
          <Links links={links} />
        </nav>
      </details>
    </>
  )
}
