// Rules shared by every TS package. Keep this the single place lint policy lives.
export const sharedRules = {
  // Reuse over copy: the API is only reached through @vernacular/api-client.
  "no-restricted-imports": [
    "error",
    {
      paths: [
        {
          name: "openapi-fetch",
          message: "Import the client from @vernacular/api-client instead.",
        },
      ],
    },
  ],
  "@typescript-eslint/consistent-type-imports": "error",
  "@typescript-eslint/no-unused-vars": ["error", { argsIgnorePattern: "^_" }],
  "import-x/no-cycle": "error",
  "import-x/no-duplicates": "error",
  complexity: ["warn", 12],
  "max-depth": ["warn", 3],
}
