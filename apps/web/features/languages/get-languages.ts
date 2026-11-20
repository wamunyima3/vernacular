import "server-only"

import { cacheLife } from "next/cache"

import { api } from "@/lib/api"

import type { Language } from "./group-by-wave"

export type LanguagesResult = { ok: true; languages: Language[] } | { ok: false; message: string }

// The list changes only when a language opens, and is the same for everyone.
// Failures throw, so they are never cached.
async function cachedLanguages(): Promise<Language[]> {
  "use cache"
  cacheLife("hours")
  const { data, error } = await api.GET("/v1/languages")
  if (error) throw new Error(error.message)
  return data.data
}

/** Fetches the language list. Returns a result instead of throwing so the UI can show the error in place. */
export async function getLanguages(): Promise<LanguagesResult> {
  try {
    return { ok: true, languages: await cachedLanguages() }
  } catch (e) {
    return {
      ok: false,
      message: e instanceof Error && e.message ? e.message : "The API could not be reached.",
    }
  }
}
