import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"

/** A card that tells the user something, with optional next steps. */
export function NoticeCard({
  title,
  children,
  actions,
  role,
}: {
  title: string
  children: React.ReactNode
  actions?: React.ReactNode
  role?: "alert"
}) {
  return (
    <Card role={role}>
      <CardHeader>
        <CardTitle className="text-balance">{title}</CardTitle>
        <CardDescription className="text-pretty">{children}</CardDescription>
      </CardHeader>
      {actions && <CardFooter className="flex flex-wrap gap-2">{actions}</CardFooter>}
    </Card>
  )
}
