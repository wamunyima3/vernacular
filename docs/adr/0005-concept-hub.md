# 0005. Concepts as the hub between languages

- Status: Accepted
- Date: 2026-11-20

## Context
Seven languages, plus English, would need 21 direct pairs or more. GATITOS is keyed on English words. Direct Zambian-to-Zambian translation is out of scope for version 1.

## Decision
Every sense and phrase links to a language-neutral `concept` (English gloss and domain). Cross-language lookup goes language → concept → language.

## Consequences
- Adding a language means linking it to concepts once, not to every other language.
- Concepts must be curated carefully: a concept that is too broad gives wrong equivalents.
- The core list of about 4,000 concepts (P0-CNT-03) becomes a key asset.
