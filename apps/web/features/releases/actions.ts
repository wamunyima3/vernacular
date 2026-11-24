"use server"

import { updateTag } from "next/cache"

import { apiErrorMessage } from "@/lib/api-error"
import { formText } from "@/lib/form-data"

import { apiForSession } from "../auth/session"
import { releasesTag } from "./get-releases"

// notes come back on failure: React resets the form after an action.
export type ReleaseState = { error?: string; notes?: string; built?: string }

/** Freezes the verified items into a new release (US-29, admin). */
export async function buildRelease(_: ReleaseState, form: FormData): Promise<ReleaseState> {
  const notes = formText(form)("notes")
  const { data, error } = await (
    await apiForSession("/admin")
  ).POST("/v1/releases", { body: { notes } })
  if (!data) return { error: apiErrorMessage(error), notes }
  updateTag(releasesTag)
  return { built: data.name }
}
