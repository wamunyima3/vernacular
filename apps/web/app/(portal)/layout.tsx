import Link from "next/link"

/** Public developer portal shell. */
export default function PortalLayout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <header className="border-b">
        <div className="mx-auto flex h-14 max-w-5xl items-center px-4">
          <Link href="/" className="font-semibold">
            Vernacular
          </Link>
        </div>
      </header>
      <main className="mx-auto w-full max-w-5xl flex-1 px-4 py-12">{children}</main>
      <footer className="border-t">
        <p className="mx-auto max-w-5xl px-4 py-6 text-sm text-pretty text-muted-foreground">
          Open source under the MIT licence. Every item says where it came from and who checked it.
        </p>
      </footer>
    </>
  )
}
