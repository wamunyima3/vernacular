import { api } from "@/lib/api"
import { proxyDownload } from "@/lib/proxy-download"

/** Streams a release's open-data zip from the API (NFR-9, US-17). Public. */
export async function GET(_: Request, { params }: RouteContext<"/releases/[name]/download">) {
  const { name } = await params
  const { response } = await api.GET("/v1/releases/{name}/download", {
    params: { path: { name } },
    parseAs: "stream",
  })
  return proxyDownload(response, "application/zip", `vernacular-${name}.zip`)
}
