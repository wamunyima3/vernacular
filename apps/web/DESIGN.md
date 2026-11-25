---
name: Vernacular
description: Human-verified words and phrases for Zambian languages, through an open API.
colors:
  primary: "oklch(0.46 0.11 155)"
  primary-foreground: "oklch(0.985 0.01 155)"
  primary-dark: "oklch(0.74 0.13 155)"
  primary-foreground-dark: "oklch(0.2 0.035 155)"
  background: "oklch(1 0 0)"
  background-dark: "oklch(0.145 0 0)"
  foreground: "oklch(0.145 0 0)"
  foreground-dark: "oklch(0.985 0 0)"
  muted: "oklch(0.97 0 0)"
  muted-foreground: "oklch(0.53 0 0)"
  border: "oklch(0.922 0 0)"
  ring: "oklch(0.55 0.12 155)"
  ring-dark: "oklch(0.62 0.12 155)"
  destructive: "oklch(0.577 0.245 27.325)"
  warning: "oklch(0.5 0.11 65)"
  warning-dark: "oklch(0.8 0.12 75)"
typography:
  headline:
    fontFamily: "Atkinson Hyperlegible Next, sans-serif"
    fontSize: "1.875rem"
    fontWeight: 600
  title:
    fontFamily: "Atkinson Hyperlegible Next, sans-serif"
    fontSize: "1.5rem"
    fontWeight: 600
  body:
    fontFamily: "Atkinson Hyperlegible Next, sans-serif"
    fontSize: "1rem"
    fontWeight: 400
    lineHeight: 1.75
  label:
    fontFamily: "Atkinson Hyperlegible Next, sans-serif"
    fontSize: "0.875rem"
    fontWeight: 500
  code:
    fontFamily: "Atkinson Hyperlegible Mono, monospace"
    fontSize: "0.875rem"
rounded:
  md: "0.5rem"
  lg: "0.625rem"
  full: "9999px"
components:
  button-primary:
    backgroundColor: "{colors.primary}"
    textColor: "{colors.primary-foreground}"
    rounded: "{rounded.lg}"
  badge-verified:
    backgroundColor: "{colors.primary}"
    textColor: "{colors.primary-foreground}"
    rounded: "{rounded.full}"
  badge-machine:
    textColor: "{colors.warning}"
    rounded: "{rounded.full}"
---

# Design

## Overview

The Field Dictionary: calm, exact and accountable, like a well-kept reference work that shows its sources. The interface stays quiet so the words, their trust status and who checked them carry the page. It must read as credible to government and NGO teams, and stay light and legible on a mid-range phone on mobile data (PRODUCT.md).

Tokens live in `app/globals.css`; components are shadcn on Base UI (`components/ui`).

## Colors

- **Deep green** (`primary`) is the only accent. It marks actions (primary buttons, focus rings, links on emphasis), progress bars and trust. Never spend it on decoration.
- Neutrals are plain grey. Secondary text is `muted-foreground`, which holds 4.5:1 on both the background and `muted` surfaces.
- **Amber** (`warning`) means machine output only. **Red** (`destructive`) means disputed or an error only.
- Dark mode follows the system setting (`prefers-color-scheme`); the dark palette is composed, not inverted: a lighter green with dark text on it.
- Every pair used for text is at least 4.5:1, and focus rings at least 3:1, in both themes.

## Typography

One family, Atkinson Hyperlegible Next, chosen because it keeps look-alike letters apart (I/l/1, rn/m): headwords in unfamiliar languages must be read exactly. Atkinson Hyperlegible Mono sets code and API output. Both load through next/font with `latin` and `latin-ext`.

Headings use `text-balance`, body text `text-pretty`, numbers in data `tabular-nums`. Docs prose is capped at `max-w-lg` (about 75 characters a line); code blocks keep the full column.

## Layout

A single centred column, `max-w-5xl` with `px-4`, under a 56px header. Docs add a 12rem left nav from `lg` up. Pages work at 390px without horizontal scroll; header links fill the bar's height as tap targets.

## Elevation & Depth

Flat. Surfaces separate by borders and the `muted` tone, not shadows.

## Shapes

Rounded rectangles at the `--radius` scale (0.625rem base) for cards, inputs and buttons; fully rounded badges.

## Components

- **Trust badge** (`features/trust/components/trust-badge.tsx`): the label always names the status, and colour is a second signal. Verified: solid green. Reviewed once: green tint. From source: grey outline. Machine: dashed amber outline. Disputed: red tint.
- Buttons, inputs, selects, checkboxes, tabs and dialogs come from `components/ui`; destructive actions confirm in an AlertDialog.

## Do's and Don'ts

- Do show where a word came from and how far it has been checked, next to the word.
- Do keep green rare enough that it still means "act" or "trusted".
- Don't add gradients, glow, decorative imagery or invented proof (partner logos, testimonials).
- Don't use colour as the only signal for a status.
- Don't add a second typeface without a role only it can perform.
