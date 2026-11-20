import "server-only"

import { headers } from "next/headers"
import { redirect } from "next/navigation"
import { cache } from "react"

import { auth } from "@/lib/auth"

/** The current session, read once per request. */
export const getSession = cache(async () => auth.api.getSession({ headers: await headers() }))

/** The current session, or a redirect to sign-in that comes back to `from`. */
export async function requireSession(from: string) {
  const session = await getSession()
  if (!session) redirect(`/sign-in?next=${encodeURIComponent(from)}`)
  return session
}
