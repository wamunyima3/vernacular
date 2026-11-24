"use server"

import { revalidatePath } from "next/cache"

import { apiErrorMessage } from "@/lib/api-error"
import { formText } from "@/lib/form-data"

import { apiForSession } from "../auth/session"

export type CreateKeyState = { error?: string; secret?: string; name?: string }

/** Creates a key and returns its secret, which the API shows only this once. */
export async function createKey(_: CreateKeyState, form: FormData): Promise<CreateKeyState> {
  const name = formText(form)("name")
  if (!name) return { error: "Give the key a name." }
  const api = await apiForSession("/account")
  const { data, error } = await api.POST("/v1/keys", { body: { name } })
  if (!data) return { error: apiErrorMessage(error), name }
  revalidatePath("/account")
  return { secret: data.secret, name: data.key.name }
}

export async function revokeKey(id: string, _: { error?: string }): Promise<{ error?: string }> {
  const api = await apiForSession("/account")
  const { error } = await api.DELETE("/v1/keys/{id}", { params: { path: { id } } })
  if (error) return { error: apiErrorMessage(error) }
  revalidatePath("/account")
  return {}
}
