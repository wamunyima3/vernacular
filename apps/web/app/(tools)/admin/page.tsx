import type { Metadata } from "next"
import Link from "next/link"
import { Suspense } from "react"

import { Button } from "@/components/ui/button"

import { AdminOverview, AdminOverviewSkeleton } from "@/features/admin/components/admin-overview"

export const metadata: Metadata = { title: "Admin" }

export default function AdminPage() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between gap-4">
        <h1 className="text-3xl font-semibold text-balance">Admin</h1>
        <Button
          variant="outline"
          size="sm"
          nativeButton={false}
          render={<Link href="/rechecks/batches" />}
        >
          Re-check batches
        </Button>
      </div>
      <Suspense fallback={<AdminOverviewSkeleton />}>
        <AdminOverview />
      </Suspense>
    </div>
  )
}
