import type { Metadata } from "next"
import { notFound } from "next/navigation"

import { docsSource } from "@/lib/docs-source"
import { mdxComponents } from "@/features/docs/mdx-components"

export function generateStaticParams() {
  return docsSource.generateParams()
}

export async function generateMetadata({
  params,
}: PageProps<"/docs/[[...slug]]">): Promise<Metadata> {
  const page = docsSource.getPage((await params).slug)
  return page ? { title: page.data.title, description: page.data.description } : {}
}

/** One MDX page from content/docs (P1-WEB-04). */
export default async function DocsPage({ params }: PageProps<"/docs/[[...slug]]">) {
  const page = docsSource.getPage((await params).slug)
  if (!page) notFound()
  const Body = page.data.body
  return (
    <article className="flex max-w-3xl flex-col gap-4">
      <header className="flex flex-col gap-3">
        <h1 className="text-3xl font-semibold text-balance">{page.data.title}</h1>
        {page.data.description && (
          <p className="text-lg text-pretty text-muted-foreground">{page.data.description}</p>
        )}
      </header>
      <Body components={mdxComponents} />
    </article>
  )
}
