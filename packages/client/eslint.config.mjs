import base from "@vernacular/eslint-config/base"

export default [
  ...base,
  // This package is the published client, so it wraps openapi-fetch itself.
  { files: ["src/**/*.ts"], rules: { "no-restricted-imports": "off" } },
  { ignores: ["dist", "src/schema.ts"] },
]
