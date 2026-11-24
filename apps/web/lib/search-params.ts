/** Search params as Next.js passes them. */
export type Params = Record<string, string | string[] | undefined>

const uuid = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i

/** A parameter's first value. */
export const param = (params: Params, k: string) =>
  Array.isArray(params[k]) ? params[k][0] : params[k]

/** A comma-separated list of IDs, dropping anything that is not one. */
export const idList = (params: Params, k: string) =>
  (param(params, k) ?? "").split(",").filter((s) => uuid.test(s))
