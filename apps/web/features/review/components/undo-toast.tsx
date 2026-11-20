"use client"

import { useEffect } from "react"
import { toast } from "sonner"

import { undo } from "../actions"
import type { LastDecision, Params } from "../review-url"

/**
 * docs/design/ui.md: no confirmation for accept or reject; Undo appears as a
 * toast for 5 seconds instead. Shown once per recorded decision.
 */
export function UndoToast({ last, params }: { last: LastDecision; params: Params }) {
  const { decision } = last
  const id = last.reviews.join(",")
  useEffect(() => {
    toast(decision === "accept" ? "Accepted" : "Rejected", {
      id, // the same decision never shows twice
      duration: 5000,
      action: {
        label: "Undo",
        onClick: async () => {
          const result = await undo(params)
          if (result?.error) toast.error(result.error)
        },
      },
    })
    // Show once per decision; params only change with it.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id])
  return null
}
