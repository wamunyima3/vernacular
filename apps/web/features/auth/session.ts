import "server-only"

import { headers } from "next/headers"
import { redirect } from "next/navigation"
import { cache } from "react"

import { api, apiAs } from "@/lib/api"
import { auth } from "@/lib/auth"

/** The current session, read once per request. */
export const getSession = cache(async () => auth.api.getSession({ headers: await headers() }))

/** The current session, or a redirect to sign-in that comes back to `from`. */
export async function requireSession(from: string) {
  const session = await getSession()
  if (!session) redirect(`/sign-in?next=${encodeURIComponent(from)}`)
  return session
}

/** An API client acting as the signed-in user, for server actions. */
export async function apiForSession(from: string) {
  const { session } = await requireSession(from)
  return apiAs(session.token)
}

/**
 * An API client for public pages: as the visitor when signed in, so their calls count
 * against their own limit, otherwise as the portal (WEB_API_KEY, ADR-0011).
 */
export async function apiForVisitor() {
  const session = await getSession()
  return session ? apiAs(session.session.token) : api
}
