"use server"

import type { Schemas } from "@vernacular/api-client"

import { apiErrorMessage } from "@/lib/api-error"
import { keptValues } from "@/lib/form-data"
import { apiAs } from "@/lib/api"

import { requireSession } from "../auth/session"
import { type ContributionErrors, contributionInput } from "./form-data"

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
