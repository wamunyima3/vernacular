import { loader } from "fumadocs-core/source"
import { docs } from "fumadocs-mdx:collections/server"

/** The written docs in content/docs, as pages and a navigation tree (ADR-0008). */
export const docsSource = loader({ baseUrl: "/docs", source: docs.toFumadocsSource() })
