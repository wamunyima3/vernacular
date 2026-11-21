"use client"

import { useActionState } from "react"

import { Button } from "@/components/ui/button"
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { LanguageField } from "@/components/language-field"
import { NativeSelect, NativeSelectOption } from "@/components/ui/native-select"
import { capitalise } from "@/lib/utils"

import { approve } from "../actions"
import { reviewerLevels } from "../levels"

/** Approve by email, since reviewers sign up themselves first. */
export function ApproveForm({ languages }: { languages: { code: string; name: string }[] }) {
  const [state, action, pending] = useActionState(approve, {})
  return (
    <form action={action}>
      <FieldGroup className="grid gap-4 sm:grid-cols-[2fr_1fr_1fr_auto] sm:items-end">
        <Field>
          <FieldLabel htmlFor="email">Email</FieldLabel>
          <Input id="email" name="email" type="email" required autoComplete="off" />
        </Field>
        <LanguageField languages={languages} />
        <Field>
          <FieldLabel htmlFor="level">Level</FieldLabel>
          <NativeSelect id="level" name="level" defaultValue="reviewer">
            {reviewerLevels.map((l) => (
              <NativeSelectOption key={l} value={l}>
                {capitalise(l)}
              </NativeSelectOption>
            ))}
          </NativeSelect>
        </Field>
        <Button type="submit" disabled={pending}>
          {pending ? "Approving…" : "Approve"}
        </Button>
      </FieldGroup>
      {state.error && <FieldError className="mt-3">{state.error}</FieldError>}
    </form>
  )
}
