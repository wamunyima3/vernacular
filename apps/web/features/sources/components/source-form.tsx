"use client"

import type { Schemas } from "@vernacular/api-client"
import { useActionState } from "react"

import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldLegend,
  FieldSet,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { NativeSelect, NativeSelectOption } from "@/components/ui/native-select"
import { Textarea } from "@/components/ui/textarea"
import { capitalise } from "@/lib/utils"

import type { SourceFormState } from "../actions"
import { allowedUses, sourceKinds } from "../form-data"

const useLabels: Record<(typeof allowedUses)[number], string> = {
  redistribute: "Redistribute",
  derive: "Make derived works",
  commercial: "Commercial use",
  train: "Train models",
}

/** Create or edit a register entry. Errors appear under the form, where the action happened. */
export function SourceForm({
  action,
  source,
  submitLabel,
}: {
  action: (state: SourceFormState, form: FormData) => Promise<SourceFormState>
  source?: Schemas["SourceInput"]
  submitLabel: string
}) {
  const [state, formAction, pending] = useActionState(action, {})
  return (
    <form action={formAction} className="flex flex-col gap-6">
      <FieldGroup>
        <Field>
          <FieldLabel htmlFor="name">Name</FieldLabel>
          <Input id="name" name="name" defaultValue={source?.name} required />
        </Field>
        <Field>
          <FieldLabel htmlFor="kind">Kind</FieldLabel>
          <NativeSelect id="kind" name="kind" defaultValue={source?.kind ?? "licensed"}>
            {sourceKinds.map((k) => (
              <NativeSelectOption key={k} value={k}>
                {capitalise(k)}
              </NativeSelectOption>
            ))}
          </NativeSelect>
        </Field>
        <Field>
          <FieldLabel htmlFor="rights_holder">Rights holder</FieldLabel>
          <Input
            id="rights_holder"
            name="rights_holder"
            defaultValue={source?.rights_holder}
            required
          />
        </Field>
        <Field>
          <FieldLabel htmlFor="licence">Licence</FieldLabel>
          <Input
            id="licence"
            name="licence"
            defaultValue={source?.licence}
            placeholder="CC BY 4.0"
            required
          />
        </Field>
        <Field>
          <FieldLabel htmlFor="url">Link</FieldLabel>
          <Input id="url" name="url" type="url" defaultValue={source?.url} />
        </Field>
        <FieldSet>
          <FieldLegend variant="label">The licence allows</FieldLegend>
          <FieldDescription>Leave all unticked until permission is signed.</FieldDescription>
          {allowedUses.map((use) => (
            <Field key={use} orientation="horizontal">
              <Checkbox
                id={`use-${use}`}
                name="allowed_uses"
                value={use}
                defaultChecked={source?.allowed_uses.includes(use)}
              />
              <FieldLabel htmlFor={`use-${use}`} className="font-normal">
                {useLabels[use]}
              </FieldLabel>
            </Field>
          ))}
        </FieldSet>
        <Field>
          <FieldLabel htmlFor="notes">Notes</FieldLabel>
          <Textarea id="notes" name="notes" defaultValue={source?.notes} />
        </Field>
      </FieldGroup>
      {state.error && <FieldError>{state.error}</FieldError>}
      {state.saved && (
        <p role="status" className="text-sm text-muted-foreground">
          Saved.
        </p>
      )}
      <Button type="submit" disabled={pending} className="self-start">
        {pending ? "Saving…" : submitLabel}
      </Button>
    </form>
  )
}
