# vernacular-ml

Training data and models for Vernacular machine translation (phase 3, AI-7).

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

## Fine-tuning

`vernacular-train` fine-tunes TranslateGemma 4B ([ADR 0013](../docs/adr/0013-translategemma.md))
with LoRA on a corpus, in both directions, then scores the dev split with chrF per language
and direction. It needs a GPU (24 GB is enough for 4B) and a Hugging Face token for an
account that has accepted the Gemma Terms of Use.

```sh
export HF_TOKEN=...
uv run --group train vernacular-train corpus/ runs/2026.1 --language bem --language nya
```

It writes the LoRA adapter and `report.json` to `runs/2026.1`. Prompts come from
`vernacular_ml/prompt.py`, which the translation service (P3-03) will use too, so
training and serving use the same prompt. The test split is kept for the benchmark (P3-05).

## Checks

```sh
make test-ml
```
