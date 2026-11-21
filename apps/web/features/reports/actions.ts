"use server"

import { redirect } from "next/navigation"

import type { NoteDecisionState } from "@/components/note-decision-form"
import { apiErrorMessage } from "@/lib/api-error"
import { formText, isOneOf } from "@/lib/form-data"

import { apiForSession } from "../auth/session"

export type DecideState = NoteDecisionState

const choices = ["confirm", "dismiss"] as const

/** Confirms or dismisses one report with the reviewer's note, then reloads the inbox. */
export async function decide(
  reportId: string,
  back: string,
  _: DecideState,
  form: FormData,
): Promise<DecideState> {
  const choice = form.get("choice")
  const note = formText(form)("note")
  if (!isOneOf(choices, choice)) return { error: "Choose whether the report has merit." }
  if (!note) return { error: "Write a note saying why." }
  const api = await apiForSession(back)
  const { error } = await api.POST("/v1/reports/{id}/decision", {
    params: { path: { id: reportId } },
    body: { confirm: choice === "confirm", note },
  })
  if (error) return { error: apiErrorMessage(error) }
  redirect(back)
}
