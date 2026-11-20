import { getSessionCookie } from "better-auth/cookies"
import { type NextRequest, NextResponse } from "next/server"

/**
 * Sends visitors without a session cookie to sign-in before any tools page renders.
 * This is only a fast first check; pages still validate the session with `requireSession`.
 */
export function proxy(request: NextRequest) {
  if (getSessionCookie(request)) return NextResponse.next()
  const url = new URL("/sign-in", request.url)
  url.searchParams.set("next", request.nextUrl.pathname + request.nextUrl.search)
  return NextResponse.redirect(url)
}

// Signed-in areas from app/(tools). Add each new tools route here.
export const config = { matcher: ["/account/:path*", "/sources/:path*"] }
