# Tonga (Citonga) style guide

**Status:** draft, written without a native speaker; needs two senior reviewers (P2-01c) before Tonga review starts. Builds on [common.md](common.md).

**Reference orthography:** Ministry of Education, _Zambian Languages: Orthography Approved by the Ministry of Education_ (1977), to be confirmed against the current edition (P0-OUT-06). This guide covers Zambian Tonga (Southern Province), not Tonga of Malawi (`tog`).

## Spelling

| Rule              | Standard                   | Also seen (store as variant) |
| ----------------- | -------------------------- | ---------------------------- |
| /tʃ/              | `c`: `cintu`               | `ch`: `chintu`               |
| Velar nasal /ŋ/   | `ng'`: `ng'ombe`           | `ŋ`, `ng`                    |
| Long vowels       | doubled: `maanzi`          | single vowel                 |
| Stop /b/ vs /β/   | `bb` for the stop: `bbuku` | `b`                          |
| Nasal + consonant | `mb`, `nd`, `nk`, `nj`     |                              |

- The language name is spelled `Citonga`; `Chitonga` (as in the languages table) is a variant.
- **Decide:** the `bb`/`b` and other doubled-consonant rules (this draft writes `ibbwe` but `mabwe`, as sources do; which is right?) (`dd`, `kk`) against the current orthography, since writers often drop them.
- **Decide:** whether `ch` should become the standard, given how people write today.

## Citation forms

- **Nouns:** the singular with its class prefix and without the augment: `muntu`, plural `bantu`. Forms with an augment (`o-`, `a-`, `i-` in some speech) are variants. Class 5 `i-` is a prefix, not an augment: `ibbwe` keeps it, because the noun has no other prefix in the singular.
- **Verbs:** the infinitive with `ku-`: `kulya` (to eat). Root: `-lya`.
- **Adjectives:** the stem with a hyphen: `-botu` (good). Agreeing forms are not separate entries.

## Noun classes

| Class | Prefix        | Example                     |
| ----- | ------------- | --------------------------- |
| 1/2   | mu- / ba-     | muntu / bantu (person)      |
| 3/4   | mu- / mi-     | musamu / misamu (medicine)  |
| 5/6   | i-, li- / ma- | ibbwe / mabwe (stone)       |
| 7/8   | ci- / zi-     | cintu / zintu (thing)       |
| 9/10  | N- / N-, zi-  | ng'ombe / ng'ombe (cow)     |
| 11    | lu-           | lwizi (river)               |
| 12/13 | ka- / tu-     | kantu / tuntu (small thing) |
| 14    | bu-           | buntu (humanity)            |
| 15    | ku-           | kulya (eating)              |

`maanzi` (water) is class 6 with no singular: record `noun_class: 6` and leave `plural` empty.

**Decide:** every example in this table, and whether class 10 plurals take `zi-` (`zing'ombe`) or stay unmarked.

## Varieties

Leave `variety` empty for the standard (Plateau) Tonga. Values: `valley` (Gwembe Valley Tonga), `toka-leya`. **Decide** with senior reviewers whether more are needed.

## Change log

| Date       | Change      | Agreed by |
| ---------- | ----------- | --------- |
| 2026-11-24 | First draft | —         |
