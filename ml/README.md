# vernacular-ml

Training data and, from P3-02b, models for Vernacular machine translation (phase 3, AI-7).

## Parallel corpus

Builds English↔language pairs from a release zip (`GET /v1/releases/{name}/download`):
each verified sense (gloss and headword), example sentence and phrase.

- Only items whose every cited source allows `train` in the source register are used.
- Pairs split into train, dev and test (90/5/5) by a hash of their English text, so a
  meaning never crosses splits and the split stays the same from release to release.
- `--exclude` takes a held-out test set, one sentence a line (FLORES for example); any
  pair whose English or text is one of those sentences, compared whole after case and
  spacing are normalised, is left out, so the test set is never trained on (DQ-10).

```sh
uv run vernacular-corpus vernacular-2026.1.zip corpus/ --exclude flores200.bem_Latn.devtest
```

It writes `corpus/<language>/{train,dev,test}.jsonl`, one `{id, language, kind, en, text}`
per line, and `corpus/manifest.json` with the counts.

## Checks

```sh
make test-ml
```
