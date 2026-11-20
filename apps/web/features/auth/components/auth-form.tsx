"use client"

import Link from "next/link"
import { useActionState } from "react"

import { Button } from "@/components/ui/button"
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"

import type { AuthFormState } from "../actions"

type Mode = "sign-in" | "sign-up"

const copy = {
  "sign-in": {
    submit: "Sign in",
    switchText: "New here?",
    switchLink: "Create an account",
    href: "/sign-up",
  },
  "sign-up": {
    submit: "Create account",
    switchText: "Already have an account?",
    switchLink: "Sign in",
    href: "/sign-in",
  },
} as const

/** Shared sign-in and sign-up form. Errors show under the form, where the action happened. */
export function AuthForm({
  mode,
  action,
  next,
}: {
  mode: Mode
  action: (state: AuthFormState, form: FormData) => Promise<AuthFormState>
  next?: string
}) {
  const [state, formAction, pending] = useActionState(action, {})
  const text = copy[mode]
  const switchHref = next ? `${text.href}?next=${encodeURIComponent(next)}` : text.href

  return (
    <form action={formAction} className="flex flex-col gap-6">
      {next && <input type="hidden" name="next" value={next} />}
      <FieldGroup>
        {mode === "sign-up" && (
          <Field>
            <FieldLabel htmlFor="name">Name</FieldLabel>
            <Input id="name" name="name" autoComplete="name" defaultValue={state.name} required />
          </Field>
        )}
        <Field>
          <FieldLabel htmlFor="email">Email</FieldLabel>
          <Input
            id="email"
            name="email"
            type="email"
            autoComplete="email"
            defaultValue={state.email}
            required
          />
        </Field>
        <Field>
          <FieldLabel htmlFor="password">Password</FieldLabel>
          <Input
            id="password"
            name="password"
            type="password"
            autoComplete={mode === "sign-in" ? "current-password" : "new-password"}
            minLength={mode === "sign-up" ? 10 : undefined}
            required
          />
        </Field>
      </FieldGroup>
      {state.error && <FieldError role="alert">{state.error}</FieldError>}
      <Button type="submit" disabled={pending}>
        {pending ? "Please wait…" : text.submit}
      </Button>
      <p className="text-sm text-pretty text-muted-foreground">
        {text.switchText}{" "}
        <Link
          href={switchHref}
          className="font-medium text-foreground underline underline-offset-4"
        >
          {text.switchLink}
        </Link>
      </p>
    </form>
  )
}
