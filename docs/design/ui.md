# UI design: portal and review tools

Task P0-ENG-07. This is the screen-level design every web task builds from. It follows the `/baseline-ui` rules (in short: shadcn/ui on Base UI, Tailwind defaults, no gradients or decorative animation, `text-balance` on headings and `text-pretty` on body, `tabular-nums` on data, skeletons for loading, errors where the action happened, `h-dvh` not `h-screen`).

## Principles

1. **Trust is always visible.** Every item shows its status. `machine` is never shown without its label, and never in the same style as `verified`.
2. **Phone first for reviewers.** Reviewers work on phones in short sessions (US-20), often on slow mobile data. The review screen fits one item on a 360 px screen with no horizontal scroll, and works on a slow 3G connection.
3. **Developer pages answer "can I use this?" in one screen:** what the data is, how good it is, and how to call it.
4. **One accent colour per view.** Neutral palette from shadcn `base-nova`; the primary colour marks the one main action.

## Shared components

| Component        | Where                        | Notes                                                                                                  |
| ---------------- | ---------------------------- | ------------------------------------------------------------------------------------------------------ |
| `SiteHeader`     | `components/site-header.tsx` | Logo left, nav right. Portal: Docs, Sign in. Tools: Review, Contribute, Account, Sign out.             |
| `TrustBadge`     | `features/trust/` (Phase 1)  | Renders `Trust.status`. Text label always, colour only as a second signal (accessible without colour). |
| `SourceList`     | `features/trust/` (Phase 1)  | Source name, licence, and model for machine items.                                                     |
| `JsonPanel`      | `features/try-it/` (Phase 1) | Raw response beside the formatted result (US-02). Monospace, copy button with `aria-label`.            |
| `LanguagePicker` | `features/languages/`        | Base UI Select. Planned languages are listed but disabled, with "planned" text.                        |
| `ErrorCard`      | `features/*/components`      | The existing pattern: card with `role="alert"`, the message, and one next step.                        |

### Trust badge

| Status     | Variant                    | Label                  |
| ---------- | -------------------------- | ---------------------- |
| `verified` | `default` (solid)          | Verified               |
| `reviewed` | `secondary`                | Reviewed once          |
| `sourced`  | `outline`                  | From source, unchecked |
| `machine`  | `outline`, dashed border   | Machine, unchecked     |
| `disputed` | `destructive` (tools only) | Disputed               |

## Portal (public)

### Home `/`

```
┌──────────────────────────────────────────────┐
│ Vernacular                     Docs  Sign in │
├──────────────────────────────────────────────┤
│ Verified words and phrases for Zambian       │
│ languages                                    │
│ One-paragraph promise. [Get an API key]      │
│                                              │
│ ┌ Try it ───────────────────────────────────┐│
│ │ [bem ▾] [ water              ] [Search]   ││
│ └───────────────────────────────────────────┘│
│                                              │
│ Languages (live verified counts)             │
│ Wave 1  [Bemba 1,204 ✓] [Nyanja 980 ✓]       │
│ Wave 2  [Tonga planned] [Lozi planned]       │
│ Wave 3  …                                    │
├──────────────────────────────────────────────┤
│ MIT · data CC BY 4.0 · contributors          │
└──────────────────────────────────────────────┘
```

### Try it `/try` (US-02, P1-WEB-02)

Desktop: two columns, formatted results left, `JsonPanel` right. Phone: tabs "Result" / "JSON".

```
[bem ▾] [ amenshi        ] [min_status: verified ▾] [Search]
┌ Result ───────────────────┐ ┌ JSON ───────────────────────┐
│ amenshi  noun · class 6   │ │ { "data": [ { "headword":   │
│ water                     │ │   "amenshi", "trust": {     │
│ [Verified] GATITOS CC BY  │ │   "status": "verified" …    │
│ matched: headword         │ │                     [Copy]  │
└───────────────────────────┘ └─────────────────────────────┘
curl / JS / Go / Python tabs with the same request
```

Empty result: "No verified entry for 'xyz' in Bemba." with one action: "Include unverified results".

### Entry `/entries/[id]` (US-16)

Headword, part of speech, noun class, plural, root, variety and register on one line. Below it, numbered senses, each with a gloss, examples and a `TrustBadge`. At the bottom: variants, sources, and a "Report a problem" button that opens a Base UI Dialog with a form (US-15).

### Docs `/docs/*` (P1-WEB-04)

Left nav on desktop, collapsible on phone. Pages: Getting started, Verification levels, Concepts, Varieties, Errors (built), Licences, API reference (Scalar, P1-WEB-03).

## Tools (signed in)

### Review `/review` (US-20, US-21)

One item per screen. The bottom action bar is fixed and respects `env(safe-area-inset-bottom)`.

```
┌──────────────────────────────┐
│ Bemba · Quick confirm   3/20 │  ← queue type and session progress (tabular-nums)
├──────────────────────────────┤
│ English                      │
│ water                        │
│                              │
│ Candidates (no provider)     │
│ ┌──────────────────────────┐ │
│ │ ○ amenshi                │ │  ← radio group (Base UI)
│ │ ○ menshi                 │ │
│ └──────────────────────────┘ │
│ + Write my own               │
│                              │
│ Part of speech  noun         │
│ Noun class      6            │
├──────────────────────────────┤
│ [Skip]  [Reject] [ Accept ]  │  ← fixed bar; Accept is the only primary button
└──────────────────────────────┘
```

- Blind: no other reviewer's decision, no provider names, never the reviewer's own items (CON-3, AI-4).
- "Write my own" opens a correction form inline, not a modal, so it works with the phone keyboard.
- Reject asks for a reason in a sheet (one tap from a fixed list, plus free text).
- No confirmation dialog for accept or reject: they are reversible by a later review. Undo appears as a toast for 5 seconds.
- Loading the next item shows a skeleton of the same card shape.
- Empty queue: "Nothing left to review in Bemba." with one action: "Switch language" or "Contribute a word".

### Contribute `/contribute` (US-18)

A single-column form: language, headword, part of speech, noun class (shown only for nouns), English gloss, an optional example and its translation, and the source. Each field's errors appear below it. When submitted, the entry gets the `sourced` badge and the page says "In the review queue."

### Disputes `/disputes` (US-23, senior reviewers)

A list of disagreements and reports, oldest first, with days open (`tabular-nums`, red after 14 days). Opening one shows each reviewer's decision side by side. Two actions resolve it: "Keep as verified" and "Send back to review"; both require a note.

### Sources `/sources` (US-25, admin)

A table with these columns: name, kind, licence, rights holder, permission document. "Add source" opens a form page, not a dialog, because the permission document is an upload. Deleting a source is destructive and uses `AlertDialog`.

### Admin `/admin` (US-26)

- **Reviewer approvals:** a table of users with a language and level per row.
- **Provider quality:** a table of acceptance and correction rates per provider and language. It is a table, not a chart, until there is enough data.

## Accessibility

- Every interactive element is reachable by keyboard, using Base UI primitives. Never hand-built.
- Icon-only buttons have `aria-label`.
- Text in the item's language carries `lang="bem"` / `lang="ny"`, so screen readers pronounce it properly.
- Colour is never the only signal: trust and day counts always have text.

## Performance budget for the review screen

- Under 100 KB of JavaScript, compressed, for the route.
- Next item prefetched while the reviewer reads the current one.
