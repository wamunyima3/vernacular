"use server"

import { redirect } from "next/navigation"

import { apiErrorMessage } from "@/lib/api-error"

import { apiForSession } from "../auth/session"

export type CreditState = { error?: string }

/** Records whether the user wants their name on the contributors page (NFR-6). */
export async function setCredit(_: CreditState, form: FormData): Promise<CreditState> {
  const api = await apiForSession("/account")
  const { error } = await api.PUT("/v1/me/credit", {
    body: { credited: form.get("credited") === "on" },
  })
  if (error) return { error: apiErrorMessage(error) }
  redirect("/account")
}
