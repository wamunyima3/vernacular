# Common rules (all languages)

## Text

- UTF-8, Unicode NFC. The normaliser enforces this; reviewers do not need to.
- Straight apostrophe `'` (U+0027) in spellings such as `ng'`. Curly apostrophes are converted.
- Tone is not marked in headwords. A tone note goes in the sense's `notes` field when it tells two words apart.
- Lower case for headwords, except proper nouns.

## Entries

| Field            | Rule                                                                                                                |
| ---------------- | ------------------------------------------------------------------------------------------------------------------- |
| `headword`       | The citation form set in the language guide (noun with prefix, verb as infinitive).                                 |
| `variants`       | Every other spelling people really write: older orthography, missing augment, `ch`/`c`.                             |
| `root`           | The stem with a leading hyphen, e.g. `-ntu`. Verbs end in the final vowel: `-lya`.                                  |
| `part_of_speech` | One of: noun, verb, adjective, adverb, pronoun, numeral, conjunction, preposition, interjection, phrase, ideophone. |
| `noun_class`     | Singular/plural class numbers in Bleek–Meinhof numbering, e.g. `1/2`.                                               |
| `plural`         | The full plural form, for nouns only.                                                                               |
| `register`       | `neutral` by default. Others: `formal`, `informal`, `town`, `respectful`, `vulgar`, `archaic`.                      |
| `variety`        | Leave empty for the standard variety. Use a value from the language guide otherwise.                                |

## Senses

- One sense per distinct meaning. If two English glosses mean the same thing, they are one sense.
- `gloss`: the shortest English equivalent(s), comma-separated, lower case, e.g. `water`.
- `definition`: optional English sentence, used when a gloss alone misleads.
- Each sense links to one concept (ADR-0005). If no concept fits, propose one; do not stretch an existing concept.

## Examples

- A full, natural sentence a native speaker would say, with an English translation.
- No personal data of real people. No hate speech, even to illustrate a vulgar word; describe the use in `notes` instead.
- Source every example: a licensed source, or "contributor" with the user ID.

## Loanwords

Loanwords in everyday use are real words and are accepted (`register: town` when they belong to urban speech). Record the source language in `notes`.

## Review decisions

- **Accept**: correct, and spelled to this guide.
- **Correct**: the meaning is right but the form is not. Fix it; the old text is kept in history.
- **Reject**: wrong meaning, not a word, or offensive without a reason to keep it.
- **Unsure**: skip. Never accept something you would not use yourself.
