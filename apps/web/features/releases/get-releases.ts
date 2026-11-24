import "server-only"

import type { Schemas } from "@vernacular/api-client"
import { cacheLife, cacheTag } from "next/cache"

import { api } from "@/lib/api"
import { settle } from "@/lib/settle"

/** Building a release refreshes this tag (features/releases/actions.ts). */
export const releasesTag = "releases"

// Releases are public and rare. Failures throw, so they are never cached.
async function cachedReleases(): Promise<Schemas["Release"][]> {
  "use cache"
  cacheLife("hours")
  cacheTag(releasesTag)
  const { data, error } = await api.GET("/v1/releases")
  if (error) throw new Error(error.message)
  return data.data
}

/** Every release, newest first, or why they could not be loaded. */
export const getReleases = () => settle(cachedReleases)
