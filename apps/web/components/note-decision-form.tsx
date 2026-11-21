"use client"

import { useActionState } from "react"

import { Button } from "@/components/ui/button"
import { Field, FieldError, FieldLabel } from "@/components/ui/field"
import { Textarea } from "@/components/ui/textarea"

export type NoteDecisionState = { error?: string }

/**
 * A decision that must say why: a note, then one submit button per choice,
 * sent as `name` (docs/design/ui.md, Disputes and Reports).
 */
export function NoteDecisionForm({
  id,
  name,
  choices,
  action: submit,
}: {
  id: string
  name: string
  choices: [{ value: string; label: string }, { value: string; label: string }]
  action: (state: NoteDecisionState, form: FormData) => Promise<NoteDecisionState>
}) {
  const [state, action, pending] = useActionState(submit, {})
  const noteId = `note-${id}`
  return (
    <form action={action} className="flex flex-col gap-3">
      <Field>
        <FieldLabel htmlFor={noteId}>Note</FieldLabel>
        <Textarea id={noteId} name="note" required maxLength={1000} />
      </Field>
      {state.error && <FieldError>{state.error}</FieldError>}
      <div className="flex flex-wrap gap-2">
        {choices.map((c, i) => (
          <Button
            key={c.value}
            type="submit"
            name={name}
            value={c.value}
            variant={i === 0 ? "default" : "outline"}
            disabled={pending}
          >
            {c.label}
          </Button>
        ))}
      </div>
    </form>
  )
}
