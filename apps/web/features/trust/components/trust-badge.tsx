import type { Schemas } from "@vernacular/api-client"

import { Badge } from "@/components/ui/badge"

type Status = Schemas["TrustStatus"]

// docs/design/ui.md: the label always says the status; colour is only a second signal.
// Green deepens with checking; amber dashed marks machine output (DESIGN.md "Trust").
const badges: Record<
  Status,
  { label: string; variant: "default" | "outline" | "destructive"; className?: string }
> = {
  verified: { label: "Verified", variant: "default" },
  reviewed: {
    label: "Reviewed once",
    variant: "outline",
    className: "border-primary/40 bg-primary/10 text-primary",
  },
  sourced: { label: "From source, unchecked", variant: "outline" },
  machine: {
    label: "Machine, unchecked",
    variant: "outline",
    className: "border-dashed border-warning text-warning",
  },
  disputed: { label: "Disputed", variant: "destructive" },
}

/** How far an item has been checked (ADR-0004). */
export function TrustBadge({ status }: { status: Status }) {
  const b = badges[status]
  return (
    <Badge variant={b.variant} className={b.className}>
      {b.label}
    </Badge>
  )
}
