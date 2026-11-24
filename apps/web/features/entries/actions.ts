"use server"

import { apiErrorMessage } from "@/lib/api-error"
import { formText } from "@/lib/form-data"

import { apiForVisitor } from "../auth/session"

export type ReportState = { error?: string; sent?: boolean }

/** Reports one meaning as wrong (US-15). Signed-out visitors report anonymously. */
export async function reportSense(
  senseId: string,
  _: ReportState,
  form: FormData,
): Promise<ReportState> {
  const reason = formText(form)("reason")
  if (reason.length < 3) return { error: "Say what is wrong, in a few words." }
  const api = await apiForVisitor()
  const { error } = await api.POST("/v1/reports", { body: { sense_id: senseId, reason } })
  return error ? { error: apiErrorMessage(error) } : { sent: true }
}
