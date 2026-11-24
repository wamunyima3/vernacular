import Link from "next/link"

import { getSession } from "../session"

/** "Account" for signed-in visitors, "Sign in" for others. */
export async function AccountLink() {
  const session = await getSession()
  return session ? (
    <Link href="/account" className="font-medium">
      Account
    </Link>
  ) : (
    <Link href="/sign-in" className="font-medium">
      Sign in
    </Link>
  )
}
