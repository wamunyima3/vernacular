"use client"

import { useActionState, useState } from "react"

import { Button } from "@/components/ui/button"
import { Field, FieldDescription, FieldError, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"

import { createKey } from "../actions"

/** Names and creates a key, then shows its secret once, with a copy button (US-01). */
export function CreateKeyForm() {
  const [state, action, pending] = useActionState(createKey, {})
  const [copied, setCopied] = useState<string>()

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
          {state.error && <FieldError>{state.error}</FieldError>}
        </Field>
      </form>
      {state.secret && (
        <Field className="rounded-lg border p-4">
          <FieldLabel htmlFor="new-key">Your new key, {state.name}</FieldLabel>
          <div className="flex gap-2">
            <Input id="new-key" readOnly value={state.secret} className="font-mono" />
            <Button
              type="button"
              variant="outline"
              onClick={() =>
                navigator.clipboard.writeText(state.secret!).then(() => setCopied(state.secret))
              }
            >
              {copied === state.secret ? "Copied" : "Copy"}
            </Button>
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
