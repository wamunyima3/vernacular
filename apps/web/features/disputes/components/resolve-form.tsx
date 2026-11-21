"use client"

import { NoteDecisionForm } from "@/components/note-decision-form"

import { resolve } from "../actions"

/** The two ways to settle a dispute; both need a note (docs/design/ui.md). */
export function ResolveForm({ senseId }: { senseId: string }) {
  return (
    <NoteDecisionForm
      id={senseId}
      name="outcome"
      choices={[
        { value: "keep", label: "Keep as verified" },
        { value: "send_back", label: "Send back to review" },
      ]}
      action={resolve.bind(null, senseId)}
    />
  )
}
