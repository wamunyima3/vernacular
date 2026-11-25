"use client"

import Link from "next/link"
import { useActionState } from "react"

import { CopyButton } from "@/components/copy-button"
import { Button } from "@/components/ui/button"
import { Field, FieldDescription, FieldError, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"

import { createKey } from "../actions"

/** Names and creates a key, then shows its secret once, with a copy button (US-01). */
export function CreateKeyForm() {
  const [state, action, pending] = useActionState(createKey, {})

  return (
    <div className="flex flex-col gap-4">
      <form action={action} className="flex flex-col gap-3">
        <Field>
          <FieldLabel htmlFor="key-name">Key name</FieldLabel>
          <div className="flex gap-2">
            <Input
              id="key-name"
              name="name"
              required
              maxLength={100}
              placeholder="My SMS app"
              defaultValue={state.secret ? "" : state.name}
              key={state.secret}
            />
            <Button type="submit" disabled={pending}>
              Create key
            </Button>
          </div>
          <FieldDescription>
            By creating a key you agree to the <Link href="/docs/terms">terms of use</Link>.
          </FieldDescription>
          {state.error && <FieldError>{state.error}</FieldError>}
        </Field>
      </form>
      {state.secret && (
        <Field className="rounded-lg border p-4">
          <FieldLabel htmlFor="new-key">Your new key, {state.name}</FieldLabel>
          <div className="flex gap-2">
            <Input id="new-key" readOnly value={state.secret} className="font-mono" />
            <CopyButton text={state.secret} label="Copy key" />
          </div>
          <FieldDescription>
            Copy it now and keep it secret. It is not shown again; if you lose it, revoke it and
            create another. Send it in the <code>X-API-Key</code> header.
          </FieldDescription>
        </Field>
      )}
    </div>
  )
}
