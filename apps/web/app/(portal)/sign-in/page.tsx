import type { Metadata } from "next"
import { Suspense } from "react"

import { AuthPage } from "@/features/auth/components/auth-page"

export const metadata: Metadata = { title: "Sign in" }

export default function SignInPage({ searchParams }: PageProps<"/sign-in">) {
  return (
    <Suspense>
      <AuthPage mode="sign-in" searchParams={searchParams} />
    </Suspense>
  )
}
