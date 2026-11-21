import type { Metadata } from "next"
import { Suspense } from "react"

import { AdminOverview, AdminOverviewSkeleton } from "@/features/admin/components/admin-overview"

export const metadata: Metadata = { title: "Admin" }

export default function AdminPage() {
  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-3xl font-semibold text-balance">Admin</h1>
      <Suspense fallback={<AdminOverviewSkeleton />}>
        <AdminOverview />
      </Suspense>
    </div>
  )
}
