import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

import { type AuthMode, authModes } from "../modes"
import { safeNext } from "../safe-next"
import { AuthForm } from "./auth-form"

/** Body of the sign-in and sign-up pages. */
export async function AuthPage({
  mode,
  searchParams,
}: {
  mode: AuthMode
  searchParams: Promise<{ next?: string | string[] }>
}) {
  const { next } = await searchParams
  const page = authModes[mode]
  return (
    <Card className="mx-auto w-full max-w-sm">
      <CardHeader>
        <CardTitle>
          <h1 className="text-xl font-semibold text-balance">{page.title}</h1>
        </CardTitle>
        <CardDescription className="text-pretty">{page.description}</CardDescription>
      </CardHeader>
      <CardContent>
        <AuthForm mode={mode} next={next ? safeNext(next) : undefined} />
      </CardContent>
    </Card>
  )
}
