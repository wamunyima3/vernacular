import { defineConfig, defineDocs } from "fumadocs-mdx/config"

/** Written docs, MDX under content/docs, served at /docs (ADR-0008). */
export const docs = defineDocs({ dir: "content/docs" })

/** Code colours that keep 4.5:1 on the muted code background (impeccable audit, WCAG 1.4.3). */
export default defineConfig({
  mdxOptions: {
    rehypeCodeOptions: { themes: { light: "github-light-high-contrast", dark: "github-dark" } },
  },
})
