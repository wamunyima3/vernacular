import type { Metadata } from "next"
import { Suspense } from "react"

import { Skeleton } from "@/components/ui/skeleton"
import { sessionWithAdmin } from "@/features/auth/admin"
import { createSource } from "@/features/sources/actions"
import { AdminOnly } from "@/features/sources/components/admin-only"
import { SourceForm } from "@/features/sources/components/source-form"

export const metadata: Metadata = { title: "Add source" }

async function NewSource() {
  const { isAdmin } = await sessionWithAdmin("/sources/new")
  if (!isAdmin) return <AdminOnly />
  return <SourceForm action={createSource} submitLabel="Add source" />
}

export default function NewSourcePage() {
  return (
    <div className="flex max-w-xl flex-col gap-6">
      <h1 className="text-2xl font-semibold text-balance">Add source</h1>
      <Suspense fallback={<Skeleton className="h-96 w-full rounded-xl" />}>
        <NewSource />
      </Suspense>
    </div>
  )
}
