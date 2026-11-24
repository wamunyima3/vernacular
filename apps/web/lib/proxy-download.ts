import "server-only"

/**
 * Passes a file the API answered with on to the browser as a download: the
 * API's error status as is, otherwise its body and caching headers, named by
 * the API's Content-Disposition or the fallback.
 */
export function proxyDownload(response: Response, type: string, filename: string) {
  if (!response.ok) return new Response(null, { status: response.status })
  const headers = new Headers({
    "Content-Type": type,
    "Content-Disposition":
      response.headers.get("Content-Disposition") ?? `attachment; filename="${filename}"`,
  })
  for (const name of ["ETag", "Cache-Control"]) {
    const value = response.headers.get(name)
    if (value) headers.set(name, value)
  }
  return new Response(response.body, { headers })
}
