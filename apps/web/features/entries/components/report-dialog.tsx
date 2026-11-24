"use client"

import { useActionState } from "react"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Field, FieldError, FieldLabel } from "@/components/ui/field"
import { Textarea } from "@/components/ui/textarea"

import { reportSense } from "../actions"

/** "Report a problem" with one meaning: a reviewer of the language reads it (US-15). */
export function ReportDialog({ senseId, meaning }: { senseId: string; meaning: string }) {
  const [state, action, pending] = useActionState(reportSense.bind(null, senseId), {})
  const fieldId = `reason-${senseId}`
  return (
    <Dialog>
      <DialogTrigger render={<Button variant="ghost" size="sm" />}>Report a problem</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Report “{meaning}”</DialogTitle>
          <DialogDescription>
            A reviewer of the language reads every report. If they agree, senior reviewers settle it
            within 14 days.
          </DialogDescription>
        </DialogHeader>
        {state.sent ? (
          <>
            <p className="text-pretty">Thank you. Your report is with the reviewers.</p>
            <DialogFooter>
              <DialogClose render={<Button />}>Close</DialogClose>
            </DialogFooter>
          </>
        ) : (
          <form action={action} className="flex flex-col gap-4">
            <Field>
              <FieldLabel htmlFor={fieldId}>What is wrong?</FieldLabel>
              <Textarea
                id={fieldId}
                name="reason"
                required
                minLength={3}
                maxLength={1000}
                placeholder="This word means moon, not water."
              />
              {state.error && <FieldError>{state.error}</FieldError>}
            </Field>
            <DialogFooter>
              <DialogClose render={<Button variant="outline" />}>Cancel</DialogClose>
              <Button type="submit" disabled={pending}>
                Send report
              </Button>
            </DialogFooter>
          </form>
        )}
      </DialogContent>
    </Dialog>
  )
}
