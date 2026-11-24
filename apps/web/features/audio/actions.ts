"use server"

import { revalidatePath } from "next/cache"

import { apiErrorMessage } from "@/lib/api-error"

import { apiForSession } from "../auth/session"

export type RecordState = { error?: string; saved?: boolean }

/** Saves the signed-in speaker's recording of an entry, replacing an earlier one. */
export async function recordClip(
  entryId: string,
  _: RecordState,
  form: FormData,
): Promise<RecordState> {
  const file = form.get("clip")
  if (!(file instanceof Blob) || file.size === 0) return { error: "Record the word first." }
  const api = await apiForSession(`/entries/${entryId}`)
  const { data, error } = await api.POST("/v1/audio", {
    params: { query: { entry_id: entryId } },
    body: file as unknown as string, // raw bytes; the contract types binary bodies as string
    bodySerializer: (body) => body,
    headers: { "Content-Type": "application/octet-stream" },
  })
  if (!data) return { error: apiErrorMessage(error) }
  revalidatePath(`/entries/${entryId}`)
  return { saved: true }
}

/** Deletes one of the speaker's clips. */
export async function deleteClip(
  entryId: string,
  id: string,
  _: { error?: string },
): Promise<{ error?: string }> {
  const api = await apiForSession(`/entries/${entryId}`)
  const { error } = await api.DELETE("/v1/audio/{id}", { params: { path: { id } } })
  if (error) return { error: apiErrorMessage(error) }
  revalidatePath(`/entries/${entryId}`)
  return {}
}

/** Withdraws the signed-in user's consent to be recorded; their clips are deleted. */
export async function withdrawConsent(_: { error?: string }): Promise<{ error?: string }> {
  const api = await apiForSession("/account")
  const { error } = await api.DELETE("/v1/me/speaker")
  if (error) return { error: apiErrorMessage(error) }
  revalidatePath("/account")
  return {}
}
