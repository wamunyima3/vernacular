/** Narrows a form value to one of a list of allowed strings. */
export const isOneOf = <T extends string>(list: readonly T[], value: unknown): value is T =>
  typeof value === "string" && (list as readonly string[]).includes(value)

/** Reads a trimmed text field, empty when missing. */
export const formText = (form: FormData) => (name: string) => String(form.get(name) ?? "").trim()

/**
 * The text a form held, to hand back from a failed server action: React
 * resets the form after an action, so the fields re-fill from these.
 */
export const keptValues = (form: FormData): Record<string, string> =>
  Object.fromEntries([...form].flatMap(([k, v]) => (typeof v === "string" ? [[k, v]] : [])))
