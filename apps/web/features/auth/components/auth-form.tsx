"use client"

import Link from "next/link"
import { useActionState } from "react"

import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"

import { signIn, signUp } from "../actions"
import { type AuthMode, authModes } from "../modes"
import { termsField } from "../terms"

const actions = { "sign-in": signIn, "sign-up": signUp } satisfies Record<AuthMode, unknown>

/** Shared sign-in and sign-up form. Errors show under the form, where the action happened. */
export function AuthForm({ mode, next }: { mode: AuthMode; next?: string }) {
  const [state, formAction, pending] = useActionState(actions[mode], {})
  const text = authModes[mode]
  const switchHref = next ? `${text.switchHref}?next=${encodeURIComponent(next)}` : text.switchHref

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
        {mode === "sign-up" && (
          <Field orientation="horizontal">
            <Checkbox id={termsField} name={termsField} defaultChecked={state.agreed} required />
            <FieldLabel htmlFor={termsField} className="font-normal">
              <span>
                I agree to the{" "}
                <Link href="/docs/terms" className="underline underline-offset-4">
                  terms of use
                </Link>
                , including the limits on machine translation.
              </span>
            </FieldLabel>
          </Field>
        )}
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
