import Link from "next/link"

import { Button } from "@/components/ui/button"
import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { apiErrorMessage } from "@/lib/api-error"
import { api } from "@/lib/api"

import { sessionWithAdmin } from "../../auth/admin"
import { AdminOnly } from "./admin-only"
import { SourceTable } from "./source-table"

export async function SourcesOverview() {
  const { isAdmin } = await sessionWithAdmin("/sources")
  if (!isAdmin) return <AdminOnly />

  const { data, error } = await api.GET("/v1/sources")
  if (!data) {
    return (
      <Card role="alert">
        <CardHeader>
          <CardTitle>The register could not be loaded</CardTitle>
          <CardDescription>{apiErrorMessage(error)} Reload the page to try again.</CardDescription>
        </CardHeader>
      </Card>
    )
  }
  return (
    <div className="flex flex-col gap-4">
      <Button render={<Link href="/sources/new" />} nativeButton={false} className="self-start">
        Add source
      </Button>
      {data.data.length === 0 ? (
        <p className="text-muted-foreground">No sources yet. Add the first licensed source.</p>
      ) : (
        <SourceTable sources={data.data} />
      )}
    </div>
  )
}
