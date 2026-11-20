import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

import { signIn, signUp } from "../actions"
import { safeNext } from "../safe-next"
import { AuthForm } from "./auth-form"

const pages = {
  "sign-in": {
    title: "Sign in",
    description: "Contribute, review and manage your API keys.",
    action: signIn,
  },
  "sign-up": {
    title: "Create an account",
    description: "Free for developers, contributors and reviewers.",
    action: signUp,
  },
} as const

/** Body of the sign-in and sign-up pages. */
export async function AuthPage({
  mode,
  searchParams,
}: {
  mode: keyof typeof pages
  searchParams: Promise<{ next?: string | string[] }>
}) {
  const { next } = await searchParams
  const page = pages[mode]
  return (
    <Card className="mx-auto w-full max-w-sm">
      <CardHeader>
        <CardTitle>
          <h1 className="text-xl font-semibold text-balance">{page.title}</h1>
        </CardTitle>
        <CardDescription className="text-pretty">{page.description}</CardDescription>
      </CardHeader>
      <CardContent>
        <AuthForm mode={mode} action={page.action} next={next ? safeNext(next) : undefined} />
      </CardContent>
    </Card>
  )
}
