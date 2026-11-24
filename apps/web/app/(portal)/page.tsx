import Link from "next/link"
import { Suspense } from "react"

import { Button } from "@/components/ui/button"

import { LanguageList } from "@/features/languages/components/language-list"
import { LanguageListSkeleton } from "@/features/languages/components/language-list-skeleton"

export default function HomePage() {
  return (
    <div className="flex flex-col gap-12">
      <section className="flex max-w-2xl flex-col gap-4">
        <h1 className="text-3xl font-semibold text-balance sm:text-4xl">
          Verified words and phrases for Zambian languages
        </h1>
        <p className="text-lg text-pretty text-muted-foreground">
          A public API for Bemba, Nyanja and the other regional languages. Nothing is marked
          verified until two native speakers have accepted it, and machine output is always
          labelled.
        </p>
        <div className="flex flex-wrap gap-3">
          <Button nativeButton={false} render={<Link href="/sign-up?next=/account" />}>
            Get an API key
          </Button>
          <Button variant="outline" nativeButton={false} render={<Link href="/credits" />}>
            See who checks the words
          </Button>
        </div>
      </section>

      <section aria-labelledby="languages-heading" className="flex flex-col gap-6">
        <h2 id="languages-heading" className="text-xl font-semibold text-balance">
          Languages
        </h2>
        <Suspense fallback={<LanguageListSkeleton />}>
          <LanguageList />
        </Suspense>
      </section>
    </div>
  )
}
