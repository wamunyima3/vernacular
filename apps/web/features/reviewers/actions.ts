"use server"

import { redirect } from "next/navigation"

import { apiErrorMessage } from "@/lib/api-error"
import { apiAs } from "@/lib/api"

import { requireSession } from "../auth/session"

export type RetrainState = { error?: string }

/** Records that a reviewer was retrained, which restarts their agreement count. */
export async function retrain(
  userId: string,
  language: string,
  _: RetrainState,
): Promise<RetrainState> {
  const { session } = await requireSession("/reviewers")
  const { error } = await apiAs(session.token).POST("/v1/reviewers/{user_id}/retraining", {
    params: { path: { user_id: userId } },
    body: { language },
  })
  if (error) return { error: apiErrorMessage(error) }
  redirect(`/reviewers?${new URLSearchParams({ language })}`)
}
