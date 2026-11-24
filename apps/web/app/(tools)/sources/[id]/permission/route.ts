import { getSession } from "@/features/auth/session"
import { apiAs } from "@/lib/api"
import { proxyDownload } from "@/lib/proxy-download"

/** Streams the permission document from the API, as the signed-in admin. */
export async function GET(_: Request, { params }: RouteContext<"/sources/[id]/permission">) {
  const session = await getSession()
  if (!session) return new Response("Sign in first", { status: 401 })
  const { id } = await params
  const { response } = await apiAs(session.session.token).GET("/v1/sources/{id}/permission", {
    params: { path: { id } },
    parseAs: "stream",
  })
  return proxyDownload(response, "application/octet-stream", `permission-${id}`)
}
