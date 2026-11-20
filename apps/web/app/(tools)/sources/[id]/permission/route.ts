import { getSession } from "@/features/auth/session"
import { apiAs } from "@/lib/api"

/** Streams the permission document from the API, as the signed-in admin. */
export async function GET(_: Request, { params }: RouteContext<"/sources/[id]/permission">) {
  const session = await getSession()
  if (!session) return new Response("Sign in first", { status: 401 })
  const { id } = await params
  const { response } = await apiAs(session.session.token).GET("/v1/sources/{id}/permission", {
    params: { path: { id } },
    parseAs: "stream",
  })
  if (!response.ok) return new Response(null, { status: response.status })
  return new Response(response.body, {
    headers: {
      "Content-Type": "application/octet-stream",
      "Content-Disposition": `attachment; filename="permission-${id}"`,
    },
  })
}
