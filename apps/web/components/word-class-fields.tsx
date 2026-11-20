"use client"

import { useState } from "react"

import { Field, FieldDescription, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { NativeSelect, NativeSelectOption } from "@/components/ui/native-select"
import { capitalise } from "@/lib/utils"
import { partsOfSpeech } from "@/lib/word-class"

/**
 * Part of speech, and the noun class when it is a noun. Read the values with
 * `wordClassInput` and the same prefix. `unsure` adds an empty first choice.
 */
export function WordClassFields({
  prefix = "",
  defaultPos = "noun",
  defaultNounClass,
  unsure = false,
}: {
  prefix?: string
  defaultPos?: string
  defaultNounClass?: string
  unsure?: boolean
}) {
  const [pos, setPos] = useState(unsure ? "" : defaultPos)
  return (
    <>
      <Field>
        <FieldLabel htmlFor={`${prefix}part_of_speech`}>Part of speech</FieldLabel>
        <NativeSelect
          id={`${prefix}part_of_speech`}
          name={`${prefix}part_of_speech`}
          value={pos}
          onChange={(e) => setPos(e.target.value)}
        >
          {unsure && <NativeSelectOption value="">Not sure</NativeSelectOption>}
          {partsOfSpeech.map((p) => (
            <NativeSelectOption key={p} value={p}>
              {capitalise(p)}
            </NativeSelectOption>
          ))}
        </NativeSelect>
      </Field>
      {pos === "noun" && (
        <Field>
          <FieldLabel htmlFor={`${prefix}noun_class`}>Noun class</FieldLabel>
          <Input
            id={`${prefix}noun_class`}
            name={`${prefix}noun_class`}
            defaultValue={defaultNounClass}
            placeholder="7/8"
            pattern="[0-9]{1,2}(/[0-9]{1,2})?"
          />
          <FieldDescription>
            Singular/plural pair, such as 1/2 or 7/8. Leave empty if unsure.
          </FieldDescription>
        </Field>
      )}
    </>
  )
}
