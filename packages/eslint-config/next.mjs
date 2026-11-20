import nextVitals from "eslint-config-next/core-web-vitals"
import nextTs from "eslint-config-next/typescript"
import importX from "eslint-plugin-import-x"

import { sharedRules } from "./rules.mjs"

/** Lint config for Next.js apps. */
export default [
  ...nextVitals,
  ...nextTs,
  {
    plugins: { "import-x": importX },
    rules: sharedRules,
  },
  { ignores: [".next/**", "out/**", "build/**", "next-env.d.ts"] },
]
