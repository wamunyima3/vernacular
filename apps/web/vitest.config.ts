import { fileURLToPath } from "node:url"

import { defineConfig } from "vitest/config"

// The same `@/` alias as tsconfig.json, so tested modules can import from lib/.
export default defineConfig({
  resolve: { alias: { "@": fileURLToPath(new URL(".", import.meta.url)) } },
})
