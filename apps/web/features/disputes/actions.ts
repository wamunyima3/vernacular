"use server"

import { redirect } from "next/navigation"

import type { NoteDecisionState } from "@/components/note-decision-form"
import { apiErrorMessage } from "@/lib/api-error"
import { formText, isOneOf } from "@/lib/form-data"

import { apiForSession } from "../auth/session"

export type ResolveState = NoteDecisionState

const outcomes = ["keep", "send_back"] as const

/** Settles one dispute with the senior reviewer's note, then reloads the list. */
export async function resolve(
  itemId: string,
  _: ResolveState,
  form: FormData,
): Promise<ResolveState> {
  const outcome = form.get("outcome")
  const note = formText(form)("note")
  if (!isOneOf(outcomes, outcome)) return { error: "Choose how to settle it." }
  if (!note) return { error: "Write a note saying why." }
  const api = await apiForSession("/disputes")
  const { error } = await api.POST("/v1/disputes/{item_id}/resolution", {
    params: { path: { item_id: itemId } },
    body: { outcome, note },
  })
  if (error) return { error: apiErrorMessage(error) }
  redirect("/disputes")
}
