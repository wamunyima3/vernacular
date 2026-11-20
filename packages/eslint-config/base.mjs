import js from "@eslint/js"
import importX from "eslint-plugin-import-x"
import tseslint from "typescript-eslint"

import { sharedRules } from "./rules.mjs"

/** Lint config for plain TS packages. */
export default tseslint.config(js.configs.recommended, ...tseslint.configs.recommended, {
  plugins: { "import-x": importX },
  rules: sharedRules,
})
