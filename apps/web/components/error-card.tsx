import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

/** A failed load, shown where the content would be, with one next step. */
export function ErrorCard({ title, message }: { title: string; message: string }) {
  return (
    <Card role="alert">
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <CardDescription className="text-pretty">
          {message} Reload the page to try again.
        </CardDescription>
      </CardHeader>
    </Card>
  )
}
