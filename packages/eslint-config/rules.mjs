import importX from "eslint-plugin-import-x"

// Lint policy shared by every TS package. Keep it in this one place.

/** Generated or vendored code that is never linted. */
export const generatedIgnores = {
  ignores: ["**/schema.d.ts", "**/components/ui/**", "**/.next/**", "**/next-env.d.ts"],
}

export const shared = {
  plugins: { "import-x": importX },
  rules: {
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
  },
}
