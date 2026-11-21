"use client"

import { useActionState } from "react"

import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Field, FieldDescription, FieldError, FieldLabel } from "@/components/ui/field"

import { setCredit } from "../actions"

/** Consent to public credit by name, which is off until the user turns it on (NFR-6). */
export function CreditConsent({ credited }: { credited: boolean }) {
  const [state, action, pending] = useActionState(setCredit, {})
  return (
    <form action={action} className="flex flex-col gap-3">
      <Field orientation="horizontal">
        <Checkbox id="credited" name="credited" defaultChecked={credited} />
        <div className="flex flex-col gap-1">
          <FieldLabel htmlFor="credited" className="font-normal">
            Show my name on the public contributors page
          </FieldLabel>
          <FieldDescription>
            Your name, with how many of your words were verified and reviews settled. Turn it off at
            any time.
          </FieldDescription>
        </div>
      </Field>
      {state.error && <FieldError>{state.error}</FieldError>}
      <Button type="submit" variant="outline" size="sm" className="self-start" disabled={pending}>
        {pending ? "Saving…" : "Save"}
      </Button>
    </form>
  )
}
