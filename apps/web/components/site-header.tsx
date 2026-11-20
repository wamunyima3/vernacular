import Link from "next/link"

/** Top bar shared by the portal and the tools. `children` fills the right-hand side. */
export function SiteHeader({ children }: { children?: React.ReactNode }) {
  return (
    <header className="border-b">
      <div className="mx-auto flex h-14 max-w-5xl items-center justify-between gap-4 px-4">
        <Link href="/" className="font-semibold">
          Vernacular
        </Link>
        <nav className="flex items-center gap-4 text-sm">{children}</nav>
      </div>
    </header>
  )
}
