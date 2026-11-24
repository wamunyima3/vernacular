"use server"

import type { Schemas } from "@vernacular/api-client"

import { apiErrorMessage } from "@/lib/api-error"
import { keptValues } from "@/lib/form-data"
import { apiAs } from "@/lib/api"

import { requireSession } from "../auth/session"
import { type ContributionErrors, contributionInput, exampleInput, licenceError } from "./form-data"

export type ContributeState = {
  error?: string
  errors?: ContributionErrors
  /** What was typed, so a failed submission keeps it (React resets forms after an action). */
  values?: Record<string, string>
  receipt?: Schemas["Contribution"] & { headword: string }
}

export async function contribute(_: ContributeState, form: FormData): Promise<ContributeState> {
  const values = keptValues(form)
  const input = contributionInput(form)
  if (!input.ok) return { errors: input.errors, values }
  const { session } = await requireSession("/contribute")
  const { data, error } = await apiAs(session.token).POST("/v1/contributions", { body: input.body })
  if (!data) return { error: apiErrorMessage(error), values }
  return { receipt: { ...data, headword: input.body.headword } }
}

export type ExampleState = { error?: string; values?: Record<string, string>; saved?: boolean }

/** Adds the contributor's example to an existing meaning, for review (LEX-5). */
export async function contributeExample(
  entryId: string,
  senseId: string,
  _: ExampleState,
  form: FormData,
): Promise<ExampleState> {
  const values = keptValues(form)
  const example = exampleInput(form)
  if (!example.text || !example.translation_en)
    return { error: "Write the sentence and its English translation.", values }
  if (form.get("agree_to_licence") !== "on") return { error: licenceError, values }
  const { session } = await requireSession(`/entries/${entryId}`)
  const { data, error } = await apiAs(session.token).POST("/v1/senses/{id}/examples", {
    params: { path: { id: senseId } },
    body: { example, agree_to_licence: true },
  })
  if (!data) return { error: apiErrorMessage(error), values }
  return { saved: true }
}
