import { docsSource } from "@/lib/docs-source"
import { DocsNav, type DocsLink } from "@/features/docs/components/docs-nav"

/** Every docs page, in the order content/docs/meta.json sets, Errors and the API reference included. */
const links: DocsLink[] = docsSource.pageTree.children.flatMap((n) =>
  n.type === "page" ? [{ href: n.url, title: String(n.name) }] : [],
)

export default function DocsLayout({ children }: LayoutProps<"/docs">) {
  return (
    <div className="grid gap-8 lg:grid-cols-[12rem_1fr]">
      <DocsNav links={links} />
      <div className="min-w-0">{children}</div>
    </div>
  )
}
