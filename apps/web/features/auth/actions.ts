"use server"

import { APIError } from "better-auth/api"
import { headers } from "next/headers"
import { redirect } from "next/navigation"

import { auth } from "@/lib/auth"

import { safeNext } from "./safe-next"

export type AuthFormState = { error?: string; email?: string; name?: string }

function field(form: FormData, name: string) {
  const value = form.get(name)
  return typeof value === "string" ? value.trim() : ""
}

/** Runs a Better Auth call and turns its errors into a message shown next to the form. */
async function attempt(state: AuthFormState, call: () => Promise<unknown>): Promise<AuthFormState> {
  try {
    await call()
    return {}
  } catch (err) {
    if (err instanceof APIError) return { ...state, error: err.body?.message ?? err.message }
    throw err
  }
}

export async function signIn(_: AuthFormState, form: FormData): Promise<AuthFormState> {
  const email = field(form, "email")
  const result = await attempt({ email }, async () =>
    auth.api.signInEmail({
      body: { email, password: form.get("password") as string },
      headers: await headers(),
    }),
  )
  if (result.error) return result
  redirect(safeNext(form.get("next")))
}

export async function signUp(_: AuthFormState, form: FormData): Promise<AuthFormState> {
  const email = field(form, "email")
  const name = field(form, "name")
  const result = await attempt({ email, name }, async () =>
    auth.api.signUpEmail({
      body: { email, name, password: form.get("password") as string },
      headers: await headers(),
    }),
  )
  if (result.error) return result
  redirect(safeNext(form.get("next")))
}

export async function signOut() {
  await auth.api.signOut({ headers: await headers() })
  redirect("/")
}
