"use client"

import Link from "next/link"
import { useActionState, useId, useRef, useState } from "react"

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { Button } from "@/components/ui/button"
import { Field, FieldDescription, FieldError, FieldLabel } from "@/components/ui/field"
import { Textarea } from "@/components/ui/textarea"

import { buildRelease } from "../actions"

/**
 * Notes, then a confirmed build: a release is public and cannot be undone (US-29).
 * The dialog renders outside the form, so its button closes it and submits the form.
 */
export function BuildReleaseForm() {
  const [state, action, pending] = useActionState(buildRelease, {})
  const [open, setOpen] = useState(false)
  const form = useRef<HTMLFormElement>(null)
  const id = useId()
  return (
    <form ref={form} action={action} className="flex max-w-xl flex-col gap-3">
      <Field>
        <FieldLabel htmlFor={`${id}-notes`}>Change log</FieldLabel>
        <Textarea id={`${id}-notes`} name="notes" maxLength={5000} defaultValue={state.notes} />
        <FieldDescription>
          What changed, in your words. Counts and changes are added for you.
        </FieldDescription>
      </Field>
      <AlertDialog open={open} onOpenChange={setOpen}>
        <AlertDialogTrigger render={<Button className="self-start" disabled={pending} />}>
          Build a release
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Build a release?</AlertDialogTitle>
            <AlertDialogDescription>
              Every verified item is frozen into a new numbered release and published. A release
              cannot be changed or removed afterwards.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                setOpen(false)
                form.current?.requestSubmit()
              }}
            >
              Build and publish
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
      {state.error && <FieldError>{state.error}</FieldError>}
      {state.built && (
        <p role="status" className="text-sm text-pretty">
          Release {state.built} is published.{" "}
          <Link href="/releases" className="underline underline-offset-4">
            See it
          </Link>
        </p>
      )}
    </form>
  )
}
