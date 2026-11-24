import { Checkbox } from "@/components/ui/checkbox"
import { Field, FieldLabel } from "@/components/ui/field"
import { Textarea } from "@/components/ui/textarea"

/**
 * An example sentence and its English (LEX-5), named `example_text` and
 * `example_translation`. idPrefix keeps ids unique when a page has several.
 */
export function ExampleFields({
  idPrefix = "",
  values = {},
  required = false,
}: {
  idPrefix?: string
  values?: { example_text?: string; example_translation?: string }
  required?: boolean
}) {
  return (
    <>
      <Field>
        <FieldLabel htmlFor={`${idPrefix}example_text`} className="font-normal">
          A sentence using the word
        </FieldLabel>
        <Textarea
          id={`${idPrefix}example_text`}
          name="example_text"
          defaultValue={values.example_text}
          maxLength={500}
          required={required}
        />
      </Field>
      <Field>
        <FieldLabel htmlFor={`${idPrefix}example_translation`} className="font-normal">
          Its English translation
        </FieldLabel>
        <Textarea
          id={`${idPrefix}example_translation`}
          name="example_translation"
          defaultValue={values.example_translation}
          maxLength={500}
          required={required}
        />
      </Field>
    </>
  )
}

/** The contributor's agreement to the dataset licence, named `agree_to_licence`. */
export function LicenceAgreement({
  idPrefix = "",
  checked,
}: {
  idPrefix?: string
  checked?: boolean
}) {
  return (
    <Field orientation="horizontal">
      <Checkbox
        id={`${idPrefix}agree_to_licence`}
        name="agree_to_licence"
        defaultChecked={checked}
      />
      <FieldLabel htmlFor={`${idPrefix}agree_to_licence`} className="font-normal">
        I agree my contribution is published under the dataset licence, with my name as contributor.
      </FieldLabel>
    </Field>
  )
}
