import { defineConfig, defineDocs } from "fumadocs-mdx/config"

/** Written docs, MDX under content/docs, served at /docs (ADR-0008). */
export const docs = defineDocs({ dir: "content/docs" })

export default defineConfig()
