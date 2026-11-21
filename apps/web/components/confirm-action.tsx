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

type State = { error?: string }

/**
 * A button for an irreversible action, confirmed in an AlertDialog
 * (baseline UI). Errors from the action show below the button.
 */
export function ConfirmAction({
  label,
  title,
  description,
  action,
  destructive = false,
  className,
}: {
  label: string
  title: string
  description: string
  action: (state: State) => Promise<State>
  destructive?: boolean
  className?: string
}) {
  const [state, formAction, pending] = useActionState(action, {})
  const variant = destructive ? "destructive" : "default"
  return (
    <div className={className}>
      <AlertDialog>
        <AlertDialogTrigger
          render={
            <Button
              variant={destructive ? "destructive" : "outline"}
              size={destructive ? "default" : "sm"}
            />
          }
        >
          {label}
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{title}</AlertDialogTitle>
            <AlertDialogDescription>{description}</AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <form action={formAction}>
              <AlertDialogAction type="submit" variant={variant} disabled={pending}>
                {label}
              </AlertDialogAction>
            </form>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
      {state.error && <FieldError>{state.error}</FieldError>}
    </div>
  )
}
