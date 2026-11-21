"use server"

import { redirect } from "next/navigation"

import { apiErrorMessage } from "@/lib/api-error"
import { apiAs } from "@/lib/api"
import { formText, isOneOf } from "@/lib/form-data"

import { requireSession } from "../auth/session"

export type ResolveState = { error?: string }

const outcomes = ["keep", "send_back"] as const

/** Settles one dispute with the senior reviewer's note, then reloads the list. */
export async function resolve(
  senseId: string,
  _: ResolveState,
  form: FormData,
): Promise<ResolveState> {
  const outcome = form.get("outcome")
  const note = formText(form)("note")
  if (!isOneOf(outcomes, outcome)) return { error: "Choose how to settle it." }
  if (!note) return { error: "Write a note saying why." }
  const { session } = await requireSession("/disputes")
  const { error } = await apiAs(session.token).POST("/v1/disputes/{sense_id}/resolution", {
    params: { path: { sense_id: senseId } },
    body: { outcome, note },
  })
  if (error) return { error: apiErrorMessage(error) }
  redirect("/disputes")
}
