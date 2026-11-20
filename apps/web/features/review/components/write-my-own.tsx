"use client"

import { useState } from "react"

import { Button } from "@/components/ui/button"
import { Field, FieldDescription, FieldLabel, FieldSet } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { WordClassFields } from "@/components/word-class-fields"

/**
 * docs/design/ui.md: "Write my own" opens a correction form inline, not in a
 * modal, so it works with the phone keyboard. Its fields belong to the card's
 * form, which posts them with `decision=correct`.
 */
export function WriteMyOwn({ pending }: { pending: boolean }) {
  const [open, setOpen] = useState(false)

  if (!open) {
    return (
      <Button
        type="button"
        variant="link"
        className="self-start px-0"
        onClick={() => setOpen(true)}
      >
        + Write my own
      </Button>
    )
  }
  return (
    <FieldSet className="rounded-lg border p-4">
      <Field>
        <FieldLabel htmlFor="correction_headword">The right word</FieldLabel>
        <Input
          id="correction_headword"
          name="correction_headword"
          maxLength={100}
          autoCapitalize="none"
          autoFocus
        />
        <FieldDescription>
          It goes to other reviewers. The candidates above stay on record.
        </FieldDescription>
      </Field>
      <WordClassFields prefix="correction_" unsure />
      <div className="flex gap-2">
        <Button type="submit" name="decision" value="correct" disabled={pending}>
          {pending ? "Sending…" : "Send correction"}
        </Button>
        <Button type="button" variant="ghost" onClick={() => setOpen(false)}>
          Cancel
        </Button>
      </div>
    </FieldSet>
  )
}
