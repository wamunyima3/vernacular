"use client"

import { useActionState, useState } from "react"

import { Button } from "@/components/ui/button"
import { FieldError, FieldSet } from "@/components/ui/field"

import { contributeExample, type ExampleState } from "../actions"
import { ExampleFields, LicenceAgreement } from "./example-fields"

/**
 * "Add an example" under a meaning on its entry page (LEX-5): opens inline, and the
 * example goes to review like any contribution.
 */
export function AddExample({ entryId, senseId }: { entryId: string; senseId: string }) {
  const [open, setOpen] = useState(false)
  const [state, action, pending] = useActionState<ExampleState, FormData>(
    contributeExample.bind(null, entryId, senseId),
    {},
  )

  if (state.saved) {
    return (
      <p role="status" className="text-sm text-pretty text-muted-foreground">
        Thank you. Your example goes to two reviewers before it is shown.
      </p>
    )
  }
  if (!open) {
    return (
      <Button variant="link" className="self-start px-0" onClick={() => setOpen(true)}>
        + Add an example
      </Button>
    )
  }
  return (
    <form action={action}>
      <FieldSet className="rounded-lg border p-4">
        <ExampleFields idPrefix={`${senseId}-`} values={state.values} required />
        <LicenceAgreement
          idPrefix={`${senseId}-`}
          checked={state.values?.agree_to_licence === "on"}
        />
        {state.error && <FieldError>{state.error}</FieldError>}
        <div className="flex gap-2">
          <Button type="submit" disabled={pending}>
            {pending ? "Sending…" : "Send for review"}
          </Button>
          <Button type="button" variant="ghost" onClick={() => setOpen(false)}>
            Cancel
          </Button>
        </div>
      </FieldSet>
    </form>
  )
}
