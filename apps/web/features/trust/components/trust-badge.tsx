import type { Schemas } from "@vernacular/api-client"

import { Badge } from "@/components/ui/badge"
import { cn } from "cn"

type Status = Schemas["TrustStatus"]

// docs/design/ui.md: the label always says the status; colour is only a second signal.
const badges: Record<
  Status,
  { label: string; variant: "default" | "secondary" | "outline" | "destructive"; dashed?: boolean }
> = {
  verified: { label: "Verified", variant: "default" },
  reviewed: { label: "Reviewed once", variant: "secondary" },
  sourced: { label: "From source, unchecked", variant: "outline" },
  machine: { label: "Machine, unchecked", variant: "outline", dashed: true },
  disputed: { label: "Disputed", variant: "destructive" },
}

/** How far an item has been checked (ADR-0004). */
export function TrustBadge({ status }: { status: Status }) {
  const b = badges[status]
  return (
    <Badge variant={b.variant} className={cn(b.dashed && "border-dashed")}>
      {b.label}
    </Badge>
  )
}
