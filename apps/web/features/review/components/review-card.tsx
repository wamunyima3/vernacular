"use client"

import type { Schemas } from "@vernacular/api-client"
import Link from "next/link"
import { useActionState } from "react"

import { Button } from "@/components/ui/button"
import { Field, FieldError, FieldLabel, FieldLegend, FieldSet } from "@/components/ui/field"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet"
import { Textarea } from "@/components/ui/textarea"

import { decide } from "../actions"
import { rejectReasons } from "../decision-input"
import { nextSession, type ReviewSession, reviewUrl } from "../review-url"

const priorityLabels: Record<Schemas["ReviewPriority"], string> = {
  quick: "Quick confirm",
  standard: "Standard review",
  full: "Full review",
}

/**
 * docs/design/ui.md, Review: one card per screen, candidates without provider
 * names, and a fixed action bar where Accept is the only primary button. One
 * form carries the card; the reject sheet's controls join it by `form` id.
 */
export function ReviewCard({
  card,
  languageName,
  session,
}: {
  card: Schemas["ReviewCard"]
  languageName: string
  session: ReviewSession
}) {
  const [state, action, pending] = useActionState(decide, {})
  const single = card.candidates.length === 1
  const submitLabel = (label: string) => (pending ? "Sending…" : label)

  return (
    <div className="flex flex-col gap-6 pb-32">
      <p className="flex justify-between gap-4 text-sm text-muted-foreground">
        <span>
          {languageName} · {priorityLabels[card.priority]}
        </span>
        <span className="tabular-nums">{session.done} reviewed</span>
      </p>

      <div className="flex flex-col gap-1">
        <p className="text-sm text-muted-foreground">English</p>
        <h1 className="text-3xl font-semibold text-balance">{card.gloss_en}</h1>
      </div>

      <form id="review" action={action}>
        <input type="hidden" name="language" value={session.language} />
        <input type="hidden" name="skip" value={session.skip.join(",")} />
        <input type="hidden" name="done" value={session.done} />
        <input type="hidden" name="concept" value={card.concept_id} />
        {card.candidates.map((c) => (
          <input key={c.sense_id} type="hidden" name="candidate" value={c.sense_id} />
        ))}
        <FieldSet>
          <FieldLegend variant="label">
            {single ? "Candidate" : "Candidates, in random order"}
          </FieldLegend>
          <RadioGroup
            name="choice"
            defaultValue={single ? card.candidates[0]?.sense_id : undefined}
          >
            {card.candidates.map((c) => (
              <FieldLabel key={c.sense_id} htmlFor={c.sense_id}>
                <Field orientation="horizontal">
                  <RadioGroupItem id={c.sense_id} value={c.sense_id} />
                  <div className="flex min-w-0 flex-col gap-0.5">
                    <span className="truncate text-lg font-medium">{c.headword}</span>
                    {(c.part_of_speech || c.noun_class) && (
                      <span className="text-sm font-normal text-muted-foreground">
                        {[c.part_of_speech, c.noun_class && `class ${c.noun_class}`]
                          .filter(Boolean)
                          .join(" · ")}
                      </span>
                    )}
                  </div>
                </Field>
              </FieldLabel>
            ))}
          </RadioGroup>
        </FieldSet>
      </form>

      <div className="fixed inset-x-0 bottom-0 z-10 border-t bg-background pb-[env(safe-area-inset-bottom)]">
        <div className="mx-auto flex max-w-5xl flex-col gap-2 px-4 py-3">
          {state.error && <FieldError>{state.error}</FieldError>}
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              nativeButton={false}
              render={<Link href={reviewUrl(nextSession(session, card.concept_id, false))} />}
            >
              Skip
            </Button>
            <Sheet>
              <SheetTrigger
                render={<Button variant="outline" className="ml-auto" disabled={pending} />}
              >
                Reject
              </SheetTrigger>
              <SheetContent side="bottom" className="pb-[env(safe-area-inset-bottom)]">
                <SheetHeader>
                  <SheetTitle>Reject {single ? "this word" : "all candidates"}</SheetTitle>
                  <SheetDescription>None of them means “{card.gloss_en}”. Why?</SheetDescription>
                </SheetHeader>
                <div className="flex flex-col gap-4 px-4">
                  <RadioGroup name="reason" form="review" aria-label="Reason">
                    {rejectReasons.map((r) => (
                      <Field key={r} orientation="horizontal">
                        <RadioGroupItem id={`reason-${r}`} value={r} />
                        <FieldLabel htmlFor={`reason-${r}`} className="font-normal">
                          {r}
                        </FieldLabel>
                      </Field>
                    ))}
                  </RadioGroup>
                  <Field>
                    <FieldLabel htmlFor="note" className="font-normal">
                      Note (optional)
                    </FieldLabel>
                    <Textarea id="note" name="note" form="review" maxLength={400} />
                  </Field>
                  {state.error && <FieldError>{state.error}</FieldError>}
                </div>
                <SheetFooter>
                  <Button
                    type="submit"
                    form="review"
                    name="decision"
                    value="reject"
                    variant="destructive"
                    disabled={pending}
                  >
                    {submitLabel("Reject")}
                  </Button>
                </SheetFooter>
              </SheetContent>
            </Sheet>
            <Button type="submit" form="review" name="decision" value="accept" disabled={pending}>
              {submitLabel("Accept")}
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
