import nextVitals from "eslint-config-next/core-web-vitals"
import nextTs from "eslint-config-next/typescript"

import { generatedIgnores, shared } from "./rules.mjs"

/** Lint config for Next.js apps. */
export default [
  ...nextVitals,
  ...nextTs,
  shared,
  generatedIgnores,
  { ignores: ["out/**", "build/**"] },
]
