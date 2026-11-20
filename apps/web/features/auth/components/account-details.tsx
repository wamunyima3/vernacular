import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { apiAs } from "@/lib/api"

import { requireSession } from "../session"

/** The signed-in user as the Go API sees them (GET /v1/me), proving the session reaches the API. */
export async function AccountDetails() {
  const { session } = await requireSession("/account")
  const { data: me, error } = await apiAs(session.token).GET("/v1/me")

  if (!me) {
    return (
      <Card role="alert">
        <CardHeader>
          <CardTitle>Your account could not be loaded</CardTitle>
          <CardDescription className="text-pretty">
            {error?.message ?? "The API is unreachable."} Reload the page to try again.
          </CardDescription>
        </CardHeader>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="truncate">{me.name}</CardTitle>
        <CardDescription className="truncate">{me.email}</CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-4">
        <dl className="grid grid-cols-[auto_1fr] gap-x-6 gap-y-2 text-sm">
          <dt className="text-muted-foreground">Role</dt>
          <dd>
            <Badge variant="outline" className="capitalize">
              {me.role}
            </Badge>
          </dd>
          <dt className="text-muted-foreground">Reviews</dt>
          <dd className="text-pretty">
            {me.approvals.length === 0
              ? "Not approved to review yet. An admin approves reviewers per language."
              : me.approvals.map((a) => `${a.language} (${a.level})`).join(", ")}
          </dd>
        </dl>
      </CardContent>
    </Card>
  )
}
