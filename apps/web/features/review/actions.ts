"use server"

import { redirect } from "next/navigation"

import { apiErrorMessage } from "@/lib/api-error"
import { apiAs } from "@/lib/api"
import { formText, keptValues } from "@/lib/form-data"

import { requireSession } from "../auth/session"
import { decisionInput } from "./decision-input"
import {
  nextSession,
  type Params,
  readLast,
  readSession,
  reviewUrl,
  undoneSession,
} from "./review-url"

/** `correcting`: the error belongs to the "Write my own" form, and shows there. */
export type DecideState = { error?: string; correcting?: boolean }

/** Records a decision on the current card, then moves to the next one. */
export async function decide(_: DecideState, form: FormData): Promise<DecideState> {
  const correcting = form.get("decision") === "correct"
  const input = decisionInput(form)
  if (!input.ok) return { error: input.error, correcting }
  const { session } = await requireSession("/review")
  const { data, error } = await apiAs(session.token).POST("/v1/reviews", { body: input.body })
  if (!data) return { error: apiErrorMessage(error), correcting }
  const { language = "", skip, done } = readSession(keptValues(form))
  const last = { reviews: data.data.map((r) => r.id), decision: input.body.decisions[0]!.decision }
  redirect(reviewUrl(nextSession({ language, skip, done }, formText(form)("concept"), true), last))
}

/** Takes back the last decision, all or none, and returns to its card. */
export async function undo(params: Params): Promise<DecideState> {
  const { language, skip, done } = readSession(params)
  const last = readLast(params)
  if (!language || !last) return { error: "Nothing to undo." }
  const { session } = await requireSession("/review")
  const { error } = await apiAs(session.token).DELETE("/v1/reviews", {
    params: { query: { ids: last.reviews } },
  })
  if (error) return { error: apiErrorMessage(error) }
  redirect(reviewUrl(undoneSession({ language, skip, done })))
}
