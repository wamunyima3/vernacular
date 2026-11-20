import js from "@eslint/js"
import tseslint from "typescript-eslint"

import { generatedIgnores, shared } from "./rules.mjs"

/** Lint config for plain TS packages. */
export default tseslint.config(
  js.configs.recommended,
  ...tseslint.configs.recommended,
  shared,
  generatedIgnores,
)
