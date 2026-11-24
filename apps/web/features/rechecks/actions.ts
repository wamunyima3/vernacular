"use server"

import type { Schemas } from "@vernacular/api-client"
import { redirect } from "next/navigation"

import type { ActionState } from "@/components/confirm-action"
import type { NoteDecisionState } from "@/components/note-decision-form"
import { apiErrorMessage } from "@/lib/api-error"
import { formText, isOneOf } from "@/lib/form-data"

import { apiForSession } from "../auth/session"
import { recheckUrl } from "./recheck-url"

const outcomes = ["pass", "fail"] as const

/** Records a senior reviewer's re-check, then shows the next one (DQ-6). */
export async function recordRecheck(
  itemId: string,
  language: string,
  _: NoteDecisionState,
  form: FormData,
): Promise<NoteDecisionState> {
  const outcome = form.get("outcome")
  const note = formText(form)("note")
  if (!isOneOf(outcomes, outcome)) return { error: "Choose whether it is still right." }
  if (outcome === "fail" && !note) return { error: "Write a note saying what is wrong." }
  const api = await apiForSession(recheckUrl(language))
  const { error } = await api.POST("/v1/rechecks/{item_id}", {
    params: { path: { item_id: itemId } },
    body: { outcome, ...(note && { note }) },
  })
  if (error) return { error: apiErrorMessage(error) }
  redirect(recheckUrl(language))
}

/** Draws this quarter's sample now, instead of waiting for the worker (admins). */
export async function drawSample(): Promise<ActionState> {
  const api = await apiForSession("/rechecks/batches")
  const { error } = await api.POST("/v1/rechecks", { body: {} })
  if (error) return { error: apiErrorMessage(error) }
  redirect("/rechecks/batches")
}

/** Sends a failing batch's verified items back to review (admins). */
export async function sendBackBatch(batch: Schemas["RecheckBatchInput"]): Promise<ActionState> {
  const api = await apiForSession("/rechecks/batches")
  const { error } = await api.POST("/v1/rechecks/batches/send-back", { body: batch })
  if (error) return { error: apiErrorMessage(error) }
  redirect(`/rechecks/batches?quarter=${batch.quarter}`)
}
