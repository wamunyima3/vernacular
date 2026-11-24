import { api } from "@/lib/api"
import { proxyDownload } from "@/lib/proxy-download"

/** A pronunciation clip, streamed from the API so the browser never needs its address. */
export async function GET(_: Request, { params }: RouteContext<"/audio/[id]">) {
  const { id } = await params
  const { response } = await api.GET("/v1/audio/{id}", {
    params: { path: { id } },
    parseAs: "stream",
  })
  return proxyDownload(response, response.headers.get("Content-Type") ?? "application/octet-stream")
}
