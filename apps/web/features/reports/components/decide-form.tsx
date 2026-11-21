"use client"

import { NoteDecisionForm } from "@/components/note-decision-form"

import { decide } from "../actions"

/** Confirm sends the meaning to the senior reviewers as a dispute; both need a note. */
export function DecideForm({ reportId, back }: { reportId: string; back: string }) {
  return (
    <NoteDecisionForm
      id={reportId}
      name="choice"
      choices={[
        { value: "confirm", label: "Confirm: send to senior reviewers" },
        { value: "dismiss", label: "Dismiss" },
      ]}
      action={decide.bind(null, reportId, back)}
    />
  )
}
