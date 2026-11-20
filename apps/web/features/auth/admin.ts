import "server-only"

import { requireSession } from "./session"

/** The session, and whether it belongs to an admin. The API enforces the role too. */
export async function sessionWithAdmin(from: string) {
  const session = await requireSession(from)
  return { session, isAdmin: session.user.role === "admin" }
}
