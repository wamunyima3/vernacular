"use server"

import type { Schemas } from "@vernacular/api-client"
import { redirect } from "next/navigation"

import { apiErrorMessage } from "@/lib/api-error"
import { apiAs } from "@/lib/api"
import { formText, keptValues } from "@/lib/form-data"

import { requireSession } from "../auth/session"

export type QualifyState = {
  error?: string
  test?: Schemas["QualificationTest"]
  result?: Schemas["QualificationResult"]
  /** The choices made, so a failed hand-in keeps them (React resets forms after an action). */
  values?: Record<string, string>
}

/** Starts a test in the chosen language. */
export async function startTest(_: QualifyState, form: FormData): Promise<QualifyState> {
  const { session } = await requireSession("/qualify")
  const { data, error } = await apiAs(session.token).POST("/v1/qualifications", {
    body: { language: formText(form)("language") },
  })
  return data ? { test: data } : { error: apiErrorMessage(error) }
}

/** Hands in the test: one chosen option per question, posted as q0, q1, … */
export async function answerTest(state: QualifyState, form: FormData): Promise<QualifyState> {
  const test = state.test
  if (!test) return state
  const text = formText(form)
  const answers = test.questions.map((_, i) => text(`q${i}`))
  const values = keptValues(form)
  if (answers.some((a) => !a))
    return { test, values, error: "Answer every question before handing in." }
  const { session } = await requireSession("/qualify")
  const { data, error } = await apiAs(session.token).POST("/v1/qualifications/{id}/answers", {
    params: { path: { id: test.id } },
    body: { answers },
  })
  return data ? { result: data } : { test, values, error: apiErrorMessage(error) }
}

/** Confirms a passed test, which approves its taker as a reviewer. */
export async function confirmTest(
  id: string,
  language: string,
  _: { error?: string },
): Promise<{ error?: string }> {
  const { session } = await requireSession("/reviewers")
  const { error } = await apiAs(session.token).POST("/v1/qualifications/{id}/confirmation", {
    params: { path: { id } },
  })
  if (error) return { error: apiErrorMessage(error) }
  redirect(`/reviewers?${new URLSearchParams({ language })}`)
}
