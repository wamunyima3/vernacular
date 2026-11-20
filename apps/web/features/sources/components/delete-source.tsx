"use client"

import { useActionState } from "react"

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
import { FieldError } from "@/components/ui/field"

import type { SourceFormState } from "../actions"

/** Deleting is irreversible, so it is confirmed in an AlertDialog. */
export function DeleteSource({
  name,
  action,
}: {
  name: string
  action: (state: SourceFormState) => Promise<SourceFormState>
}) {
  const [state, formAction, pending] = useActionState(action, {})
  return (
    <div className="flex flex-col gap-3">
      <AlertDialog>
        <AlertDialogTrigger render={<Button variant="destructive" className="self-start" />}>
          Delete source
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete {name}?</AlertDialogTitle>
            <AlertDialogDescription>
              This removes it from the register. A source that items cite cannot be deleted.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <form action={formAction}>
              <AlertDialogAction type="submit" variant="destructive" disabled={pending}>
                Delete
              </AlertDialogAction>
            </form>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
      {state.error && <FieldError>{state.error}</FieldError>}
    </div>
  )
}
