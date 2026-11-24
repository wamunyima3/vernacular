import type { Metadata } from "next"
import { Suspense } from "react"

import { PerLanguage } from "@/components/per-language"
import {
  MissingExamples,
  MissingExamplesSkeleton,
} from "@/features/examples/components/missing-examples"

export const metadata: Metadata = {
  title: "Needs examples",
  description: "Core meanings that still need a verified example sentence, per language.",
}

export default function ExamplesPage(props: PageProps<"/examples">) {
  return (
    <div className="mx-auto flex max-w-2xl flex-col gap-6">
      <div className="flex flex-col gap-2">
        <h1 className="text-3xl font-semibold text-balance">Needs examples</h1>
        <p className="text-pretty text-muted-foreground">
          Every core meaning should have a sentence showing the word in use, with its English. Open
          a word and choose “Add an example”; two reviewers check it before it is shown.
        </p>
      </div>
      <Suspense fallback={<MissingExamplesSkeleton />}>
        <PerLanguage
          searchParams={props.searchParams}
          path="/examples"
          fallback={<MissingExamplesSkeleton />}
          render={(language) => <MissingExamples language={language} />}
        />
      </Suspense>
    </div>
  )
}
