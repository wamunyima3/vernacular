import type { Metadata } from "next"
import Link from "next/link"
import { notFound } from "next/navigation"
import { Suspense } from "react"

import { ErrorCard } from "@/components/error-card"
import { SignInNotice } from "@/components/sign-in-notice"
import { NoticeCard } from "@/components/notice-card"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { apiErrorMessage } from "@/lib/api-error"
import { parseMinStatus } from "@/lib/min-status"
import { apiForVisitor } from "@/features/auth/session"
import { EntryView } from "@/features/entries/components/entry-view"
import { getLanguages, languageName } from "@/features/languages/get-languages"

export const metadata: Metadata = { title: "Entry" }

async function Entry({ params, searchParams }: PageProps<"/entries/[id]">) {
  const [{ id }, { min_status }, api] = await Promise.all([params, searchParams, apiForVisitor()])
  const minStatus = parseMinStatus(min_status)
  const get = (min_status: typeof minStatus) =>
    api.GET("/v1/entries/{id}", { params: { path: { id }, query: { min_status } } })
  const [{ data, error, response }, languages] = await Promise.all([get(minStatus), getLanguages()])
  if (data) return <EntryView entry={data} languageName={languageName(languages, data.language)} />
  // 404 means no such entry, or none of its meanings checked enough: ask at the widest level.
  if (response.status === 404 && minStatus !== "machine" && (await get("machine")).data) {
    return (
      <NoticeCard
        title="No verified meanings yet"
        actions={
          <Button nativeButton={false} render={<Link href={`/entries/${id}?min_status=machine`} />}>
            Show unverified meanings
          </Button>
        }
      >
        This word has meanings that two reviewers have not both accepted yet.
      </NoticeCard>
    )
  }
  if (response.status === 404 || response.status === 400) notFound()
  if (response.status === 401) {
    return (
      <SignInNotice title="Sign in to read entries" next={`/entries/${id}`}>
        Entries come from the live API, which needs an account or an API key.
      </SignInNotice>
    )
  }
  return <ErrorCard title="The entry could not be loaded" message={apiErrorMessage(error)} />
}

export default function EntryPage(props: PageProps<"/entries/[id]">) {
  return (
    <Suspense fallback={<Skeleton className="h-96 w-full max-w-3xl" />}>
      <Entry {...props} />
    </Suspense>
  )
}
