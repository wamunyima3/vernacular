import "server-only"

import type { Schemas } from "@vernacular/api-client"
import { cacheLife } from "next/cache"

import { api } from "@/lib/api"
import { settle } from "@/lib/settle"

// Credit changes slowly and is the same for everyone. Failures throw, so they are never cached.
async function cachedCredits(language: string): Promise<Schemas["Credit"][]> {
  "use cache"
  cacheLife("hours")
  const { data, error } = await api.GET("/v1/credits", { params: { query: { language } } })
  if (error) throw new Error(error.message)
  return data.credits
}

/** A language's credit list, or why it could not be loaded. */
export const getCredits = (language: string) => settle(() => cachedCredits(language))
