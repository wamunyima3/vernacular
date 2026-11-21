import { Field, FieldLabel } from "@/components/ui/field"
import { NativeSelect, NativeSelectOption } from "@/components/ui/native-select"

/** A required language choice, posted as `language`. */
export function LanguageField({
  languages,
  defaultValue,
}: {
  languages: { code: string; name: string }[]
  defaultValue?: string
}) {
  return (
    <Field>
      <FieldLabel htmlFor="language">Language</FieldLabel>
      <NativeSelect id="language" name="language" required defaultValue={defaultValue}>
        {languages.map((l) => (
          <NativeSelectOption key={l.code} value={l.code}>
            {l.name}
          </NativeSelectOption>
        ))}
      </NativeSelect>
    </Field>
  )
}
