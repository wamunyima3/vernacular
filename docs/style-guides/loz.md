# Lozi (Silozi) style guide

**Status:** draft, written without a native speaker; needs two senior reviewers (P2-01c) before Lozi review starts. Builds on [common.md](common.md).

**Reference orthography:** Ministry of Education, _Zambian Languages: Orthography Approved by the Ministry of Education_ (1977), to be confirmed against the current edition (P0-OUT-06). Lozi grew from Sotho-Tswana (Kololo) and Luyana, so its spelling and noun classes differ from the other Zambian languages here.

## Spelling

| Rule            | Standard    | Also seen (store as variant) |
| --------------- | ----------- | ---------------------------- |
| /tʃ/            | `c`: `kuca` | `ch`: `kucha`                |
| Velar nasal /ŋ/ | `ñ`: `ñaka` | `ng'`: `ng'aka`, `ŋ`         |
| Long vowels     | not written |                              |

- **Decide:** the examples in this table, and whether `ñ` is still the norm in print and on phones, where it is hard to type.
- **Decide:** whether `ch` should become the standard, given how people write today.
- Search folds accents away, so `ñaka` is also found as `naka`, and could meet a word spelled with a plain `n`. The variant `ng'aka` is found as written.

## Citation forms

- **Nouns:** the singular with its class prefix: `mutu`, plural `batu`. Lozi has no augment.
- **Verbs:** the infinitive with `ku-`: `kuca` (to eat). Root: `-ca`.
- **Adjectives:** the stem with a hyphen. Agreeing forms are not separate entries.

## Noun classes

| Class | Prefix     | Example                |
| ----- | ---------- | ---------------------- |
| 1/2   | mu- / ba-  | mutu / batu (person)   |
| 3/4   | mu- / mi-  |                        |
| 5/6   | li- / ma-  | litapi / matapi (fish) |
| 7/8   | si- / li-  | sika / lika (thing)    |
| 9/10  | none / li- | komu / likomu (cow)    |
| 11    | lu-        |                        |
| 14    | bu-        | bupilo (life)          |
| 15    | ku-        | kuca (to eat)          |

`mezi` (water) is class 6 with no singular: record `noun_class: 6` and leave `plural` empty.

Class 8 and class 10 plurals both take `li-`, like class 5 singulars, so `li-` alone does not tell a reviewer the class. **Decide:** every example in this table, and examples for the empty rows (3/4, 11).

## Varieties

Leave `variety` empty for standard Lozi. **Decide** with senior reviewers whether any varieties need values.

## Change log

| Date       | Change      | Agreed by |
| ---------- | ----------- | --------- |
| 2026-11-24 | First draft | —         |
