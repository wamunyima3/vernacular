"use client"

import type { Schemas } from "@vernacular/api-client"
import Link from "next/link"
import { useActionState } from "react"

import { cn } from "cn"

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
import { candidateDetail } from "@/lib/word-class"

import { decide } from "../actions"
import { rejectReasons } from "../decision-input"
import { kinds } from "../kinds"
import { nextSession, type ReviewSession, reviewUrl } from "../review-url"
import { WriteMyOwn } from "./write-my-own"

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
  const { noun, heading, wraps, rejectWhy } = kinds[card.kind]
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
        <p className="text-sm text-muted-foreground">{heading}</p>
        <h1 className="text-3xl font-semibold text-balance">{card.gloss_en}</h1>
      </div>

      <form id="review" action={action} className="flex flex-col gap-4">
        <input type="hidden" name="language" value={session.language} />
        <input type="hidden" name="skip" value={session.skip.join(",")} />
        <input type="hidden" name="done" value={session.done} />
        <input type="hidden" name="concept" value={card.concept_id} />
        <input type="hidden" name="kind" value={card.kind} />
        {card.candidates.map((c) => (
          <input key={c.item_id} type="hidden" name="candidate" value={c.item_id} />
        ))}
        <FieldSet>
          <FieldLegend variant="label">
            {single ? "Candidate" : "Candidates, in random order"}
          </FieldLegend>
          <RadioGroup name="choice" defaultValue={single ? card.candidates[0]?.item_id : undefined}>
            {card.candidates.map((c) => (
              <FieldLabel key={c.item_id} htmlFor={c.item_id}>
                <Field orientation="horizontal">
                  <RadioGroupItem id={c.item_id} value={c.item_id} />
                  <div className="flex min-w-0 flex-col gap-0.5">
                    <span className={cn("text-lg font-medium", wraps ? "text-pretty" : "truncate")}>
                      {c.headword}
                    </span>
                    {candidateDetail(c) && (
                      <span className="text-sm font-normal text-pretty text-muted-foreground">
                        {candidateDetail(c)}
                      </span>
                    )}
                  </div>
                </Field>
              </FieldLabel>
            ))}
          </RadioGroup>
        </FieldSet>
        <WriteMyOwn
          pending={pending}
          kind={card.kind}
          error={state.correcting ? state.error : undefined}
        />
      </form>

      <div className="fixed inset-x-0 bottom-0 z-10 border-t bg-background pb-[env(safe-area-inset-bottom)]">
        <div className="mx-auto flex max-w-5xl flex-col gap-2 px-4 py-3">
          {state.error && !state.correcting && <FieldError>{state.error}</FieldError>}
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
                  <SheetTitle>Reject {single ? `this ${noun}` : "all candidates"}</SheetTitle>
                  <SheetDescription>{rejectWhy(card.gloss_en)}</SheetDescription>
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
