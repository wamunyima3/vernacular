"""The one prompt the fine-tuned model is trained and served with (ADR 0013).

TranslateGemma's raw prompt names the languages in a ``<translate>`` tag. Its chat
template wants ISO 639-1 codes, which Bemba, Tonga and Lozi do not have, so we use
our own codes here: ``en`` for English and the API's ISO 639-3 codes for the rest.
Training builds prompts with this module, and the translation service (P3-03) will
too, so the two cannot drift apart.
"""

from __future__ import annotations

ENGLISH = "eng"

#: Ends the model's answer.
END = "<end_of_turn>"


def tag_code(language: str) -> str:
    """The code a language goes by in the prompt."""
    return "en" if language == ENGLISH else language


def prompt(text: str, source: str, target: str) -> str:
    """The model's input for translating text from source to target (API codes)."""
    return (
        f"<start_of_turn>user\n<translate source_lang={tag_code(source)} target_lang={tag_code(target)}>\n"
        f"{text}\n</translate><end_of_turn>\n<start_of_turn>model\n"
    )
