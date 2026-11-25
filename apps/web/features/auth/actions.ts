"use server"

import { APIError } from "better-auth/api"
import { headers } from "next/headers"
import { redirect } from "next/navigation"

import { auth } from "@/lib/auth"

import { safeNext } from "./safe-next"
import { termsError, termsField } from "./terms"

export type AuthFormState = { error?: string; email?: string; name?: string; agreed?: boolean }

function field(form: FormData, name: string) {
  const value = form.get(name)
  return typeof value === "string" ? value.trim() : ""
}

/**
 * Runs a Better Auth call. On success it redirects to the form's `next` path; on failure
 * it returns the message to show next to the form, with the typed values kept.
 */
async function attempt(
  form: FormData,
  kept: AuthFormState,
  call: (requestHeaders: Headers) => Promise<unknown>,
): Promise<AuthFormState> {
  try {
    await call(await headers())
  } catch (err) {
    if (err instanceof APIError) return { ...kept, error: err.body?.message ?? err.message }
    throw err
  }
  redirect(safeNext(form.get("next")))
}

export async function signIn(_: AuthFormState, form: FormData): Promise<AuthFormState> {
  const email = field(form, "email")
  return attempt(form, { email }, (h) =>
    auth.api.signInEmail({ body: { email, password: String(form.get("password")) }, headers: h }),
  )
}

export async function signUp(_: AuthFormState, form: FormData): Promise<AuthFormState> {
  const email = field(form, "email")
  const name = field(form, "name")
  if (form.get(termsField) !== "on") return { email, name, error: termsError }
  return attempt(form, { email, name, agreed: true }, (h) =>
    auth.api.signUpEmail({
      body: { email, name, password: String(form.get("password")) },
      headers: h,
    }),
  )
}

export async function signOut() {
  await auth.api.signOut({ headers: await headers() })
  redirect("/")
}
