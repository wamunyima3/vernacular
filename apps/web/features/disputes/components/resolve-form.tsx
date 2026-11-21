"use client"

import { useActionState } from "react"

import { Button } from "@/components/ui/button"
import { Field, FieldError, FieldLabel } from "@/components/ui/field"
import { Textarea } from "@/components/ui/textarea"

import { resolve } from "../actions"

/** The two ways to settle a dispute; both need a note (docs/design/ui.md). */
export function ResolveForm({ senseId }: { senseId: string }) {
  const [state, action, pending] = useActionState(resolve.bind(null, senseId), {})
  const noteId = `note-${senseId}`
  return (
    <form action={action} className="flex flex-col gap-3">
      <Field>
        <FieldLabel htmlFor={noteId}>Note</FieldLabel>
        <Textarea id={noteId} name="note" required maxLength={1000} />
      </Field>
      {state.error && <FieldError>{state.error}</FieldError>}
      <div className="flex flex-wrap gap-2">
        <Button type="submit" name="outcome" value="keep" disabled={pending}>
          Keep as verified
        </Button>
        <Button type="submit" name="outcome" value="send_back" variant="outline" disabled={pending}>
          Send back to review
        </Button>
      </div>
    </form>
  )
}
