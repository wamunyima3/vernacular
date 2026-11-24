import type { Schemas } from "@vernacular/api-client"

import { cn } from "cn"

import { minStatuses } from "@/lib/min-status"
import { LanguageField } from "@/components/language-field"
import { Button } from "@/components/ui/button"
import { Field, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { NativeSelect, NativeSelectOption } from "@/components/ui/native-select"

export type TryQuery = { language: string; q: string; minStatus: Schemas["MinStatus"] }

/** The URL query of a search, as the form sends it and the API reads it. */
export const searchParamsOf = (q: TryQuery) =>
  new URLSearchParams({ language: q.language, q: q.q, min_status: q.minStatus })

/**
 * The search form of "Try it" (US-02). A plain GET form to /try, so a search is a link
 * that can be shared. `compact` leaves out the verification choice, for the home page.
 */
export function TryForm({
  languages,
  query,
  compact = false,
}: {
  languages: { code: string; name: string }[]
  query?: Partial<TryQuery>
  compact?: boolean
}) {
  return (
    <form
      action="/try"
      className={cn(
        "grid gap-3 sm:items-end",
        compact ? "sm:grid-cols-[10rem_1fr_auto]" : "sm:grid-cols-[10rem_1fr_14rem_auto]",
      )}
    >
      <LanguageField languages={languages} defaultValue={query?.language} />
      <Field>
        <FieldLabel htmlFor="q">Word</FieldLabel>
        <Input
          id="q"
          name="q"
          type="search"
          required
          maxLength={100}
          placeholder="amenshi"
          defaultValue={query?.q}
          autoComplete="off"
        />
      </Field>
      {!compact && (
        <Field>
          <FieldLabel htmlFor="min_status">Show</FieldLabel>
          <NativeSelect id="min_status" name="min_status" defaultValue={query?.minStatus}>
            {minStatuses.map((s) => (
              <NativeSelectOption key={s.value} value={s.value}>
                {s.label}
              </NativeSelectOption>
            ))}
          </NativeSelect>
        </Field>
      )}
      <Button type="submit">Search</Button>
    </form>
  )
}
