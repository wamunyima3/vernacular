"use client"

import type { Schemas } from "@vernacular/api-client"
import { useActionState, useState } from "react"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
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
import { WordClassFields } from "@/components/word-class-fields"

import { TrustBadge } from "../../trust/components/trust-badge"
import { type ContributeState, contribute } from "../actions"

/** docs/design/ui.md, Contribute: one column, errors under their field. */
export function ContributeForm({ languages }: { languages: Schemas["Language"][] }) {
  const [state, action, pending] = useActionState(contribute, {})
  const v = state.values ?? {}
  // Each submission returns a new state object, so "Contribute another"
  // hides this receipt without hiding the next one.
  const [dismissed, setDismissed] = useState<ContributeState>()

  if (state.receipt && state !== dismissed) {
    const r = state.receipt
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex flex-wrap items-center gap-2 text-balance">
            {r.headword} <TrustBadge status={r.trust.status} />
          </CardTitle>
          <CardDescription className="text-pretty">
            {r.new_sense
              ? "In the review queue. Two reviewers will check it before it is verified."
              : "This word and meaning were already known. Your contribution counts as agreement."}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Button onClick={() => setDismissed(state)}>Contribute another</Button>
        </CardContent>
      </Card>
    )
  }

  return (
    <form action={action} className="flex flex-col gap-6">
      <FieldGroup>
        <Field>
          <FieldLabel htmlFor="language">Language</FieldLabel>
          <NativeSelect id="language" name="language" required defaultValue={v.language}>
            {languages.map((l) => (
              <NativeSelectOption key={l.code} value={l.code}>
                {l.name}
              </NativeSelectOption>
            ))}
          </NativeSelect>
        </Field>
        <Field>
          <FieldLabel htmlFor="headword">Word</FieldLabel>
          <Input
            id="headword"
            name="headword"
            defaultValue={v.headword}
            required
            maxLength={100}
            autoCapitalize="none"
          />
          <FieldDescription>The dictionary form, spelled as the style guide says.</FieldDescription>
        </Field>
        <WordClassFields defaultNounClass={v.noun_class} />
        <Field>
          <FieldLabel htmlFor="gloss_en">English meaning</FieldLabel>
          <Input id="gloss_en" name="gloss_en" defaultValue={v.gloss_en} required maxLength={200} />
        </Field>
        <FieldSet>
          <FieldLegend variant="label">Example (optional)</FieldLegend>
          <Field>
            <FieldLabel htmlFor="example_text" className="font-normal">
              A sentence using the word
            </FieldLabel>
            <Textarea
              id="example_text"
              name="example_text"
              defaultValue={v.example_text}
              maxLength={500}
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="example_translation" className="font-normal">
              Its English translation
            </FieldLabel>
            <Textarea
              id="example_translation"
              name="example_translation"
              defaultValue={v.example_translation}
              maxLength={500}
            />
          </Field>
          {state.errors?.example && <FieldError>{state.errors.example}</FieldError>}
        </FieldSet>
        <Field orientation="horizontal">
          <Checkbox
            id="agree_to_licence"
            name="agree_to_licence"
            defaultChecked={v.agree_to_licence === "on"}
          />
          <FieldLabel htmlFor="agree_to_licence" className="font-normal">
            I agree my contribution is published under the dataset licence, with my name as
            contributor.
          </FieldLabel>
        </Field>
        {state.errors?.agree_to_licence && <FieldError>{state.errors.agree_to_licence}</FieldError>}
      </FieldGroup>
      {state.error && <FieldError>{state.error}</FieldError>}
      <Button type="submit" disabled={pending} className="self-start">
        {pending ? "Sending…" : "Send for review"}
      </Button>
    </form>
  )
}
