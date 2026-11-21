"use server"

import { redirect } from "next/navigation"

import { apiErrorMessage } from "@/lib/api-error"
import { formText, isOneOf } from "@/lib/form-data"

import { apiForSession } from "../auth/session"
import { reviewerLevels } from "./levels"

export type AdminState = { error?: string }

/** Approves a signed-up user to review a language, or changes their level (CON-2). */
export async function approve(_: AdminState, form: FormData): Promise<AdminState> {
  const text = formText(form)
  const level = form.get("level")
  const email = text("email")
  const language = text("language")
  if (!email || !language || !isOneOf(reviewerLevels, level)) {
    return { error: "Give an email, a language and a level." }
  }
  const { error } = await (
    await apiForSession("/admin")
  ).PUT("/v1/approvals", { body: { email, language, level } })
  if (error) return { error: apiErrorMessage(error) }
  redirect("/admin")
}

/** Removes a reviewer approval. */
export async function revoke(userId: string, language: string, _: AdminState): Promise<AdminState> {
  const { error } = await (
    await apiForSession("/admin")
  ).DELETE("/v1/approvals/{user_id}/{language}", {
    params: { path: { user_id: userId, language } },
  })
  if (error) return { error: apiErrorMessage(error) }
  redirect("/admin")
}
