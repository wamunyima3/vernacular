import "server-only"

import { api } from "@/lib/api"

import type { Language } from "./group-by-wave"

export type LanguagesResult = { ok: true; languages: Language[] } | { ok: false; message: string }

/** Fetches the language list. Returns a result instead of throwing so the UI can show the error in place. */
export async function getLanguages(): Promise<LanguagesResult> {
  try {
    const { data, error } = await api.GET("/v1/languages")
    if (error) return { ok: false, message: error.message }
    return { ok: true, languages: data.data }
  } catch {
    return { ok: false, message: "The API could not be reached." }
  }
}
