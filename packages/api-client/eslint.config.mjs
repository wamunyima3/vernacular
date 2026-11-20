import base from "@vernacular/eslint-config/base"

export default [
  ...base,
  { ignores: ["src/schema.d.ts"] },
  // This package is the one place allowed to wrap openapi-fetch.
  { rules: { "no-restricted-imports": "off" } },
]
