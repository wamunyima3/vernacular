"use client"

import { useActionState } from "react"

import { Button } from "@/components/ui/button"
import { Field, FieldDescription, FieldError, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"

import type { SourceFormState } from "../actions"

/** Upload or replace the signed permission document. */
export function PermissionUpload({
  action,
  hasDocument,
  downloadHref,
}: {
  action: (state: SourceFormState, form: FormData) => Promise<SourceFormState>
  hasDocument: boolean
  downloadHref: string
}) {
  const [state, formAction, pending] = useActionState(action, {})
  return (
    <form action={formAction} className="flex flex-col gap-4">
      <Field>
        <FieldLabel htmlFor="file">Signed permission document</FieldLabel>
        <Input id="file" name="file" type="file" accept="application/pdf,image/png,image/jpeg" />
        <FieldDescription>
          PDF, PNG or JPEG, up to 10 MB.{" "}
          {hasDocument ? (
            <a href={downloadHref} className="underline underline-offset-4">
              Download the current file
            </a>
          ) : (
            "Nothing on file yet."
          )}
        </FieldDescription>
      </Field>
      {state.error && <FieldError>{state.error}</FieldError>}
      {state.saved && (
        <p role="status" className="text-sm text-muted-foreground">
          Uploaded.
        </p>
      )}
      <Button type="submit" variant="outline" disabled={pending} className="self-start">
        {pending ? "Uploading…" : hasDocument ? "Replace file" : "Upload file"}
      </Button>
    </form>
  )
}
