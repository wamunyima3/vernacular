import type { Schemas } from "@vernacular/api-client"

/** Global account roles, typed against the contract so a new role fails the type check here. */
export type Role = Schemas["Role"]

export const roles = {
  developer: "developer",
  contributor: "contributor",
  admin: "admin",
} as const satisfies Record<Role, Role>
