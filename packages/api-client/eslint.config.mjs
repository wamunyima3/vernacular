import base from "@vernacular/eslint-config/base"

export default [
  ...base,
  // The client wrapper is the one file allowed to import openapi-fetch.
  { files: ["src/index.ts"], rules: { "no-restricted-imports": "off" } },
]
