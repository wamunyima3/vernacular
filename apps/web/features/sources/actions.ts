"use server"

import { revalidatePath } from "next/cache"
import { redirect } from "next/navigation"

import { apiErrorMessage } from "@/lib/api-error"
import { apiAs } from "@/lib/api"

import { requireSession } from "../auth/session"
import { sourceInput } from "./form-data"

export type SourceFormState = { error?: string; saved?: boolean }

async function client() {
  const { session } = await requireSession("/sources")
  return apiAs(session.token)
}

export async function createSource(_: SourceFormState, form: FormData): Promise<SourceFormState> {
  const { data, error } = await (await client()).POST("/v1/sources", { body: sourceInput(form) })
  if (!data) return { error: apiErrorMessage(error) }
  revalidatePath("/sources")
  redirect(`/sources/${data.id}`)
}

export async function updateSource(
  id: string,
  _: SourceFormState,
  form: FormData,
): Promise<SourceFormState> {
  const { data, error } = await (
    await client()
  ).PUT("/v1/sources/{id}", {
    params: { path: { id } },
    body: sourceInput(form),
  })
  if (!data) return { error: apiErrorMessage(error) }
  revalidatePath("/sources")
  return { saved: true }
}

export async function uploadPermission(
  id: string,
  _: SourceFormState,
  form: FormData,
): Promise<SourceFormState> {
  const file = form.get("file")
  if (!(file instanceof File) || file.size === 0) return { error: "Choose a file to upload." }
  const { data, error } = await (
    await client()
  ).PUT("/v1/sources/{id}/permission", {
    params: { path: { id } },
    body: file as unknown as string, // raw bytes; the contract types binary bodies as string
    bodySerializer: (body) => body,
    headers: { "Content-Type": "application/octet-stream" },
  })
  if (!data) return { error: apiErrorMessage(error) }
  revalidatePath(`/sources/${id}`)
  return { saved: true }
}

export async function deleteSource(id: string, _: SourceFormState): Promise<SourceFormState> {
  const { error, response } = await (
    await client()
  ).DELETE("/v1/sources/{id}", {
    params: { path: { id } },
  })
  if (!response.ok) return { error: apiErrorMessage(error) }
  revalidatePath("/sources")
  redirect("/sources")
}
