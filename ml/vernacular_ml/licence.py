"""The files the Gemma Terms of Use require beside a fine-tuned adapter (ADR 0013).

§3.1(d) wants a NOTICE file with every copy, and §3.1(c) a prominent notice that
the files were modified. §3.1(a) and (b) pass the use restrictions and the terms
on to whoever receives the adapter, so the model card links both.
"""

from __future__ import annotations

from collections.abc import Sequence
from pathlib import Path

from .prompt import ENGLISH, LANGUAGES

TERMS = "https://ai.google.dev/gemma/terms"
PROHIBITED_USE = "https://ai.google.dev/gemma/prohibited_use_policy"

NOTICE = f"Gemma is provided under and subject to the Gemma Terms of Use found at {TERMS.removeprefix('https://')}\n"


def model_card(base: str, languages: Sequence[str]) -> str:
    """A Hugging Face model card for the adapter, stating what it is modified from."""
    tags = "\n".join(f"- {code}" for code in [LANGUAGES[ENGLISH][1], *languages])
    names = ", ".join(LANGUAGES[code][0] for code in languages)
    return f"""---
base_model: {base}
library_name: peft
license: gemma
language:
{tags}
---

# Vernacular machine translation adapter

**Modified from [{base}](https://huggingface.co/{base}).** This LoRA adapter was
fine-tuned by the Vernacular project on human-verified English pairs in
{names}. It is not the original Gemma model and is not endorsed by Google.

## Licence

Gemma is provided under and subject to the [Gemma Terms of Use]({TERMS}). This
adapter is a Model Derivative under those terms, so they apply to it and to
anyone you share it with: give recipients a copy of the terms and keep these
notices.

You must not use this adapter, or anything built from it, for the restricted uses
in the [Gemma Prohibited Use Policy]({PROHIBITED_USE}), or in breach of applicable
law. These restrictions bind every user and every further distribution.

The training data is the Vernacular dataset (CC BY 4.0) and the training code is
MIT; neither changes the terms that govern these weights.

## Use

The prompt the adapter was trained with is built by `vernacular_ml.prompt`; serve
it with `vernacular_ml.service`. Its output is machine translation and is labelled
as such by the Vernacular API.
"""


def write(out: Path, base: str, languages: Sequence[str]) -> None:
    """Write NOTICE and the model card (README.md) into the adapter directory,
    replacing the card PEFT generates."""
    (out / "NOTICE").write_text(NOTICE, encoding="utf-8")
    (out / "README.md").write_text(model_card(base, languages), encoding="utf-8")
