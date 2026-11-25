import Link from "next/link"

/** Top bar shared by the portal and the tools. `children` fills the right-hand side. */
export function SiteHeader({ children }: { children?: React.ReactNode }) {
  return (
    <header className="border-b">
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:absolute focus:top-2 focus:z-50 focus:left-2 focus:rounded-md focus:bg-background focus:px-3 focus:py-2 focus:ring-2 focus:ring-ring"
      >
        Skip to content
      </a>
      <div className="mx-auto flex h-14 max-w-5xl items-center justify-between gap-4 px-4">
        <Link href="/" className="font-semibold">
          Vernacular
        </Link>
        {/* On narrow phones the links scroll inside the bar instead of widening the page. */}
        <nav className="-mr-4 flex min-w-0 items-center gap-3 self-stretch overflow-x-auto pr-4 text-sm whitespace-nowrap sm:gap-4 [&>a]:flex [&>a]:h-full [&>a]:items-center [&>a]:focus-visible:-outline-offset-2">
          {children}
        </nav>
      </div>
    </header>
  )
}
