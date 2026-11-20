"use server"

import { redirect } from "next/navigation"

import { apiErrorMessage } from "@/lib/api-error"
import { apiAs } from "@/lib/api"
import { formText, keptValues } from "@/lib/form-data"

import { requireSession } from "../auth/session"
import { decisionInput } from "./decision-input"
import { nextSession, readSession, reviewUrl } from "./review-url"

export type DecideState = { error?: string }

/** Records a decision on the current card, then moves to the next one. */
export async function decide(_: DecideState, form: FormData): Promise<DecideState> {
  const input = decisionInput(form)
  if (!input.ok) return { error: input.error }
  const { session } = await requireSession("/review")
  const { data, error } = await apiAs(session.token).POST("/v1/reviews", { body: input.body })
  if (!data) return { error: apiErrorMessage(error) }
  const { language = "", ...rest } = readSession(keptValues(form))
  redirect(reviewUrl(nextSession({ language, ...rest }, formText(form)("concept"), true)))
}
