import { join } from "node:path"

import { createMDX } from "fumadocs-mdx/next"
import type { NextConfig } from "next"

const nextConfig: NextConfig = {
  // A self-contained server for the container image (apps/web/Dockerfile), traced from the
  // monorepo root so workspace packages are included.
  output: "standalone",
  outputFileTracingRoot: join(import.meta.dirname, "../.."),
  cacheComponents: true,
  partialPrefetching: true,
  experimental: {
    // Permission documents (10 MB) go through a server action and the proxy.
    serverActions: { bodySizeLimit: "11mb" },
    proxyClientMaxBodySize: "11mb",
  },
  turbopack: {
    rules: {
      "*.css": {
        loaders: ["@tailwindcss/turbopack"],
        as: "*.css",
      },
    },
  },
}

// MDX docs under content/docs (ADR-0008).
export default createMDX()(nextConfig)
