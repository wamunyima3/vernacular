import Link from "next/link"
import { Suspense } from "react"

import { AccountLink } from "@/features/auth/components/account-link"
import { SiteHeader } from "@/components/site-header"

/** Public developer portal shell. */
export default function PortalLayout({ children }: LayoutProps<"/">) {
  return (
    <>
      <SiteHeader>
        <Link href="/try" className="text-muted-foreground hover:text-foreground">
          Try it
        </Link>
        <Link href="/docs/api" className="text-muted-foreground hover:text-foreground">
          API reference
        </Link>
        <Link href="/docs/errors" className="text-muted-foreground hover:text-foreground">
          Docs
        </Link>
        <Suspense fallback={<span className="w-14" />}>
          <AccountLink />
        </Suspense>
      </SiteHeader>
      <main className="mx-auto w-full max-w-5xl flex-1 px-4 py-12">{children}</main>
      <footer className="border-t">
        <p className="mx-auto max-w-5xl px-4 py-6 text-sm text-pretty text-muted-foreground">
          Open source under the MIT licence. Every item says where it came from and who checked it.{" "}
          <Link href="/credits" className="underline underline-offset-4 hover:text-foreground">
            Contributors
          </Link>
        </p>
      </footer>
    </>
  )
}
