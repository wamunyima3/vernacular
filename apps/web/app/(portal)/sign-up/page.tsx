import type { Metadata } from "next"
import { Suspense } from "react"

import { AuthPage } from "@/features/auth/components/auth-page"

export const metadata: Metadata = { title: "Create an account" }

export default function SignUpPage({ searchParams }: PageProps<"/sign-up">) {
  return (
    <Suspense>
      <AuthPage mode="sign-up" searchParams={searchParams} />
    </Suspense>
  )
}
