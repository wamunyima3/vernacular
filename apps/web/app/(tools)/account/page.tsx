import type { Metadata } from "next"
import { Suspense } from "react"

import { Skeleton } from "@/components/ui/skeleton"
import { AccountDetails } from "@/features/auth/components/account-details"

export const metadata: Metadata = { title: "Account" }

export default function AccountPage() {
  return (
    <div className="flex max-w-xl flex-col gap-6">
      <h1 className="text-2xl font-semibold text-balance">Account</h1>
      <Suspense fallback={<Skeleton className="h-40 w-full rounded-xl" />}>
        <AccountDetails />
      </Suspense>
    </div>
  )
}
