import type { MDXComponents } from "mdx/types"
import Link from "next/link"

import { cn } from "cn"

import { ExampleTabs } from "./components/example-tabs"

/** Docs typography in the app's own styles, so the docs share one design system (ADR-0008). */
export const mdxComponents: MDXComponents = {
  h2: (p) => <h2 className="mt-6 scroll-mt-6 text-2xl font-semibold text-balance" {...p} />,
  h3: (p) => <h3 className="mt-4 scroll-mt-6 text-lg font-semibold text-balance" {...p} />,
  p: (p) => <p className="leading-7 text-pretty" {...p} />,
  ul: (p) => <ul className="ml-6 flex list-disc flex-col gap-1" {...p} />,
  ol: (p) => <ol className="ml-6 flex list-decimal flex-col gap-1" {...p} />,
  li: (p) => <li className="pl-1 text-pretty" {...p} />,
  a: ({ href = "", ...p }) => (
    <Link href={href} className="underline underline-offset-4 hover:text-muted-foreground" {...p} />
  ),
  code: ({ className, ...p }) => (
    <code className={cn("rounded bg-muted px-1 py-0.5 font-mono text-[0.9em]", className)} {...p} />
  ),
  // Shiki's highlighted blocks (globals.css colours them); `icon` is meant for Fumadocs UI.
  pre: ({
    className,
    icon: _icon,
    style: _style,
    ...p
  }: React.ComponentProps<"pre"> & { icon?: string }) => (
    <pre
      className={cn(
        className,
        "overflow-x-auto rounded-lg bg-muted p-4 font-mono text-sm [&_code]:bg-transparent [&_code]:p-0",
      )}
      {...p}
    />
  ),
  table: (p) => (
    <div className="overflow-x-auto">
      <table className="w-full border-collapse text-sm" {...p} />
    </div>
  ),
  th: (p) => <th className="border-b px-3 py-2 text-left font-medium" {...p} />,
  td: (p) => <td className="border-b px-3 py-2 align-top text-pretty" {...p} />,
  ExampleTabs,
  blockquote: (p) => <blockquote className="border-l-2 pl-4 text-muted-foreground" {...p} />,
}
