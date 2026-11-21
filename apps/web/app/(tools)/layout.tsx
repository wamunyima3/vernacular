import Link from "next/link"

import { SiteHeader } from "@/components/site-header"
import { Button } from "@/components/ui/button"
import { signOut } from "@/features/auth/actions"

/** Shell for signed-in tools. Each page checks the session itself with `requireSession`. */
export default function ToolsLayout({ children }: LayoutProps<"/">) {
  return (
    <>
      <SiteHeader>
        <Link href="/review" className="text-muted-foreground hover:text-foreground">
          Review
        </Link>
        <Link href="/reports" className="text-muted-foreground hover:text-foreground">
          Reports
        </Link>
        <Link href="/disputes" className="text-muted-foreground hover:text-foreground">
          Disputes
        </Link>
        <Link href="/contribute" className="text-muted-foreground hover:text-foreground">
          Contribute
        </Link>
        <Link href="/sources" className="text-muted-foreground hover:text-foreground">
          Sources
        </Link>
        <Link href="/account" className="text-muted-foreground hover:text-foreground">
          Account
        </Link>
        <form action={signOut}>
          <Button type="submit" variant="outline" size="sm">
            Sign out
          </Button>
        </form>
      </SiteHeader>
      <main className="mx-auto w-full max-w-5xl flex-1 px-4 py-12">{children}</main>
    </>
  )
}
