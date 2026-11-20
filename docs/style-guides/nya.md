# Nyanja (Cinyanja) style guide

**Status:** draft, needs two senior reviewers (P0-OUT-07). Builds on [common.md](common.md).

**Reference orthography:** Ministry of Education, _Zambian Languages: Orthography Approved by the Ministry of Education_ (1977), to be confirmed against the current edition (P0-OUT-06). Zambian Nyanja is closely related to Malawian Chichewa, but spelling differs; this guide follows the Zambian orthography.

## Spelling

| Rule                | Standard (Zambia) | Also seen (store as variant)                     |
| ------------------- | ----------------- | ------------------------------------------------ |
| /tʃ/                | `c`: `cinthu`     | `ch`: `chinthu` (common, and standard in Malawi) |
| /w/ vs Malawian /β/ | `w`: `awiri`      | `ŵ`: `aŵiri` (Malawi)                            |
| Velar nasal /ŋ/     | `ng'`: `ng'ombe`  | `ŋ`                                              |
| Aspiration          | `ph`, `th`, `kh`  |                                                  |
| Affricates          | `dz`, `ts`        |                                                  |

- Malawian spellings found in sources (`ŵ`, `ch`) are kept as variants so search finds them, and the headword follows Zambian spelling.
- **Decide:** whether `ch` should become the standard, given how people write today.

## Citation forms

- **Nouns:** the singular with its prefix: `munthu`, plural `anthu`. Nyanja has no augment.
- **Verbs:** the infinitive with `ku-`: `kudya` (to eat). Root: `-dya`.
- **Adjectives:** the stem with a hyphen: `-bwino`, `-kulu`.

## Noun classes

| Class | Prefix      | Example                       |
| ----- | ----------- | ----------------------------- |
| 1/2   | mu- / a-    | munthu / anthu (person)       |
| 3/4   | mu- / mi-   | mtengo / mitengo (tree)       |
| 5/6   | (li-) / ma- | dzina / maina (name)          |
| 7/8   | ci- / zi-   | cinthu / zinthu (thing)       |
| 9/10  | N- / N-     | nyumba / nyumba (house)       |
| 12/13 | ka- / ti-   | kamwana / tiana (small child) |
| 14    | u-          | ufa (flour)                   |
| 15    | ku-         | kudya (eating)                |

`madzi` (water) is class 6 with no singular: record `noun_class: 6` and leave `plural` empty.

`mtengo` / `mutengo`: class 3 prefixes are written `m-` before most consonants. **Decide** whether `mu-` spellings are variants or errors.

## Varieties

Leave `variety` empty for standard Nyanja. Values: `town` for Lusaka urban Nyanja (heavy English and Bemba borrowing, e.g. `ma-` plurals on English words), `chewa` for rural Eastern Province Chewa usage. Malawian Chichewa forms are variants, not a variety.

## Change log

| Date       | Change      | Agreed by |
| ---------- | ----------- | --------- |
| 2026-11-20 | First draft | —         |
