import Link from "next/link"

/** Top bar shared by the portal and the tools. `children` fills the right-hand side. */
export function SiteHeader({ children }: { children?: React.ReactNode }) {
  return (
    <header className="border-b">
      <div className="mx-auto flex h-14 max-w-5xl items-center justify-between gap-4 px-4">
        <Link href="/" className="font-semibold">
          Vernacular
        </Link>
        {/* On narrow phones the links scroll inside the bar instead of widening the page. */}
        <nav className="-mr-4 flex min-w-0 items-center gap-3 overflow-x-auto pr-4 text-sm whitespace-nowrap sm:gap-4">
          {children}
        </nav>
      </div>
    </header>
  )
}
