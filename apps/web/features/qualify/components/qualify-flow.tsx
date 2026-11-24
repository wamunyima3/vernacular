"use client"

import type { Schemas } from "@vernacular/api-client"
import Link from "next/link"
import { useActionState } from "react"

import { LanguageField } from "@/components/language-field"
import { NoticeCard } from "@/components/notice-card"
import { Button } from "@/components/ui/button"
import { Field, FieldError, FieldLabel, FieldLegend, FieldSet } from "@/components/ui/field"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"

import { type QualifyState, answerTest, startTest } from "../actions"

/**
 * The reviewer test (DQ-3, US-24): choose a language, pick the word for
 * each meaning, see the score. One action state carries the whole flow.
 */
export function QualifyFlow({ languages }: { languages: { code: string; name: string }[] }) {
  const [state, action, pending] = useActionState(
    (s: QualifyState, form: FormData) => (s.test ? answerTest(s, form) : startTest(s, form)),
    {},
  )
  if (state.result) return <Score result={state.result} />
  if (!state.test) {
    return (
      <form action={action} className="flex flex-col gap-6">
        <LanguageField languages={languages} />
        {state.error && <FieldError>{state.error}</FieldError>}
        <Button type="submit" disabled={pending} className="self-start">
          {pending ? "Starting…" : "Start the test"}
        </Button>
      </form>
    )
  }
  const { test } = state
  const name = languages.find((l) => l.code === test.language)?.name ?? test.language
  return (
    <form action={action} className="flex flex-col gap-8">
      <p className="text-pretty text-muted-foreground">
        For each meaning, pick the {name} word. You hand the test in once.
      </p>
      <ol className="flex flex-col gap-8">
        {test.questions.map((q, i) => (
          <li key={i}>
            <FieldSet>
              <FieldLegend className="text-pretty">
                <span className="text-muted-foreground tabular-nums">{i + 1}.</span> {q.gloss_en}
              </FieldLegend>
              <RadioGroup
                name={`q${i}`}
                defaultValue={state.values?.[`q${i}`]}
                className="grid gap-2 sm:grid-cols-2"
              >
                {q.options.map((o) => (
                  <FieldLabel key={o.id} htmlFor={`q${i}-${o.id}`}>
                    <Field orientation="horizontal">
                      <RadioGroupItem id={`q${i}-${o.id}`} value={o.id} />
                      <span lang={test.language} className="truncate font-medium">
                        {o.headword}
                      </span>
                    </Field>
                  </FieldLabel>
                ))}
              </RadioGroup>
            </FieldSet>
          </li>
        ))}
      </ol>
      {state.error && <FieldError>{state.error}</FieldError>}
      <Button type="submit" disabled={pending} className="self-start">
        {pending ? "Handing in…" : "Hand in"}
      </Button>
    </form>
  )
}

function Score({ result }: { result: Schemas["QualificationResult"] }) {
  const score = (
    <span className="tabular-nums">
      {result.score} of {result.total}
    </span>
  )
  const contribute = (
    <Button variant="outline" nativeButton={false} render={<Link href="/contribute" />}>
      Contribute a word
    </Button>
  )
  return result.passed ? (
    <NoticeCard title="You passed" actions={contribute}>
      You got {score} right. A senior reviewer of the language confirms your pass; then the review
      queue opens for you.
    </NoticeCard>
  ) : (
    <NoticeCard title="Not this time" actions={contribute}>
      You got {score} right, and a pass needs {result.pass_mark}. You can take the test again
      tomorrow.
    </NoticeCard>
  )
}
