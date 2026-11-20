# Bemba (Icibemba) style guide

**Status:** draft, needs two senior reviewers (P0-OUT-07). Builds on [common.md](common.md).

**Reference orthography:** Ministry of Education, _Zambian Languages: Orthography Approved by the Ministry of Education_ (1977), to be confirmed against the current edition (P0-OUT-06).

## Spelling

| Rule              | Standard                     | Also seen (store as variant) |
| ----------------- | ---------------------------- | ---------------------------- |
| /tʃ/              | `c`: `icintu`                | `ch`: `ichintu`              |
| Velar nasal /ŋ/   | `ng'`: `ng'ombe`             | `ŋ`, `ng`                    |
| Long vowels       | doubled: `ukuseeka`          | single vowel                 |
| /β/, /w/          | `b`, `w`                     |                              |
| Nasal + consonant | `mb`, `nd`, `nk`, `ns`, `nj` |                              |

- GATITOS spells the language name `Ichibemba`. Our headword is `icibemba`, with `ichibemba` as a variant. Popular writing on phones uses `ch` widely, so search must match both (LEX-4).
- **Decide:** whether `ch` should become the standard, given how people write today.

## Citation forms

- **Nouns:** with the augment (initial vowel) and class prefix, as spoken alone: `umuntu`, plural `abantu`. The form without the augment (`muntu`) is a variant.
- **Verbs:** the infinitive with `uku-`: `ukulya` (to eat). Root: `-lya`.
- **Adjectives:** the stem with a hyphen: `-suma` (good). Agreeing forms are not separate entries.

## Noun classes

| Class | Prefix (with augment) | Example                       |
| ----- | --------------------- | ----------------------------- |
| 1/2   | umu- / aba-           | umuntu / abantu (person)      |
| 3/4   | umu- / imi-           | umuti / imiti (tree)          |
| 5/6   | i- / ama-             | ilinso / amenso (eye)         |
| 7/8   | ici- / ifi-           | icintu / ifintu (thing)       |
| 9/10  | in- / in-             | ing'ombe / ing'ombe (cow)     |
| 11    | ulu-                  | ulupi (palm of hand)          |
| 12/13 | aka- / utu-           | akantu / utuntu (small thing) |
| 14    | ubu-                  | ubuntu (humanity)             |
| 15    | uku-                  | ukulya (eating)               |

`amenshi` (water) is class 6 with no singular: record `noun_class: 6` and leave `plural` empty.

## Varieties

Leave `variety` empty for standard Bemba. Values: `copperbelt` (urban Copperbelt usage with English and Swahili loans), `luapula`, `northern`. **Decide** with senior reviewers whether more are needed.

## Change log

| Date       | Change      | Agreed by |
| ---------- | ----------- | --------- |
| 2026-11-20 | First draft | —         |
