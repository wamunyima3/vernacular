// No "server-only" import: the Better Auth CLI loads this file to generate the schema.
// The pg import already keeps it out of client bundles.
import { betterAuth } from "better-auth"
import { nextCookies } from "better-auth/next-js"
import { admin } from "better-auth/plugins"
import { Pool } from "pg"

import { roles } from "@/features/auth/roles"

/** Maps Better Auth's camelCase field names to our snake_case columns (ADR-0007). */
function snake<const F extends string>(...fields: F[]) {
  return Object.fromEntries(
    fields.map((f) => [f, f.replace(/[A-Z]/g, (c) => `_${c.toLowerCase()}`)]),
  ) as Record<F, string>
}

const timestamps = ["createdAt", "updatedAt"] as const

/** Server-side auth. Better Auth owns the users, sessions, accounts and verifications tables. */
export const auth = betterAuth({
  appName: "Vernacular",
  database: new Pool({ connectionString: process.env.DATABASE_URL }),
  emailAndPassword: { enabled: true, minPasswordLength: 10 },
  user: { modelName: "users", fields: snake("emailVerified", ...timestamps) },
  session: {
    modelName: "sessions",
    fields: snake("expiresAt", "ipAddress", "userAgent", "userId", ...timestamps),
  },
  account: {
    modelName: "accounts",
    fields: snake(
      "accountId",
      "providerId",
      "userId",
      "accessToken",
      "refreshToken",
      "idToken",
      "accessTokenExpiresAt",
      "refreshTokenExpiresAt",
      ...timestamps,
    ),
  },
  verification: {
    modelName: "verifications",
    fields: snake("expiresAt", ...timestamps),
  },
  plugins: [
    admin({
      defaultRole: roles.developer,
      adminRoles: [roles.admin],
      schema: {
        user: { fields: snake("banReason", "banExpires") },
        session: { fields: snake("impersonatedBy") },
      },
    }),
    nextCookies(), // must be last: lets server actions set the session cookie
  ],
})
