import type { Schemas } from "@vernacular/api-client"

export const sourceKinds = [
  "licensed",
  "contributor",
  "machine",
] as const satisfies readonly Schemas["SourceKind"][]
export const allowedUses = [
  "redistribute",
  "derive",
  "commercial",
  "train",
] as const satisfies readonly Schemas["AllowedUse"][]

const isOneOf = <T extends string>(list: readonly T[], value: unknown): value is T =>
  typeof value === "string" && (list as readonly string[]).includes(value)

/** Reads the source form into the contract's input type. The API validates it again. */
export function sourceInput(form: FormData): Schemas["SourceInput"] {
  const text = (name: string) => String(form.get(name) ?? "").trim()
  const kind = form.get("kind")
  const url = text("url")
  return {
    kind: isOneOf(sourceKinds, kind) ? kind : "licensed",
    name: text("name"),
    rights_holder: text("rights_holder"),
    licence: text("licence"),
    ...(url && { url }),
    allowed_uses: form.getAll("allowed_uses").filter((u) => isOneOf(allowedUses, u)),
    notes: text("notes"),
  }
}
