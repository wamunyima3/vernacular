import type { Metadata } from "next"
import { Suspense } from "react"

import { Skeleton } from "@/components/ui/skeleton"
import { AccountDetails } from "@/features/auth/components/account-details"
import { CreateKeyForm } from "@/features/keys/components/create-key-form"
import { KeyList, KeyListSkeleton } from "@/features/keys/components/key-list"

export const metadata: Metadata = { title: "Account" }

export default function AccountPage() {
  return (
    <div className="flex max-w-xl flex-col gap-6">
      <h1 className="text-2xl font-semibold text-balance">Account</h1>
      <Suspense fallback={<Skeleton className="h-40 w-full rounded-xl" />}>
        <AccountDetails />
      </Suspense>
      <section className="flex flex-col gap-4" aria-labelledby="keys">
        <div className="flex flex-col gap-1">
          <h2 id="keys" className="text-xl font-semibold text-balance">
            API keys
          </h2>
          <p className="text-sm text-pretty text-muted-foreground">
            Programs call the API with a key. Keep one key per app, so you can revoke one without
            breaking the others.
          </p>
        </div>
        <CreateKeyForm />
        <Suspense fallback={<KeyListSkeleton />}>
          <KeyList />
        </Suspense>
      </section>
    </div>
  )
}
