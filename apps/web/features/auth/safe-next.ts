/** Where to go after signing in. Only same-site paths, so the link cannot send people elsewhere. */
export function safeNext(next: unknown, fallback = "/account"): string {
  return typeof next === "string" && next.startsWith("/") && !next.startsWith("//")
    ? next
    : fallback
}
