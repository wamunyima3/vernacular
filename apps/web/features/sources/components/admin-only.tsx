import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

/** Shown to signed-in users who are not admins. */
export function AdminOnly() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Admins only</CardTitle>
        <CardDescription className="text-pretty">
          Only admins manage the source register. Ask an admin if a source is missing.
        </CardDescription>
      </CardHeader>
    </Card>
  )
}
