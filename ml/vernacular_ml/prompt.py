"""The one prompt the fine-tuned model is trained and served with (ADR 0013).

It is the text TranslateGemma's chat template (``chat_template.jinja`` in
google/translategemma-4b-it) produces, built here by hand because the template
only accepts the languages in its own table, and Bemba, Tonga and Lozi are not
there (its ``to`` is Tongan). Training and the translation service both build
prompts with this module, so the two cannot drift apart.
"""

from __future__ import annotations

ENGLISH = "eng"

#: Ends the model's answer.
END = "<end_of_turn>"

#: Each language's name and the code the prompt calls it by: ISO 639-1 for
#: English, as the template has it, and the API's ISO 639-3 codes for the rest.
LANGUAGES = {
    ENGLISH: ("English", "en"),
    "bem": ("Bemba", "bem"),
    "nya": ("Nyanja", "nya"),
    "toi": ("Tonga", "toi"),
    "loz": ("Lozi", "loz"),
}


def _named(language: str) -> tuple[str, str]:
    try:
        return LANGUAGES[language]
    except KeyError:
        raise ValueError(f"no prompt name for language {language!r}; add it to LANGUAGES") from None


def prompt(text: str, source: str, target: str) -> str:
    """The model's input for translating text from source to target (API codes),
    starting with the <bos> token the template puts first."""
    src, src_code = _named(source)
    tgt, tgt_code = _named(target)
    return (
        "<bos><start_of_turn>user\n"
        f"You are a professional {src} ({src_code}) to {tgt} ({tgt_code}) translator. "
        f"Your goal is to accurately convey the meaning and nuances of the original {src} text "
        f"while adhering to {tgt} grammar, vocabulary, and cultural sensitivities.\n"
        f"Produce only the {tgt} translation, without any additional explanations or commentary. "
        f"Please translate the following {src} text into {tgt}:\n\n\n{text.strip()}<end_of_turn>\n"
        "<start_of_turn>model\n"
    )


def answer(generated: str) -> str:
    """The model's translation, without anything after its end-of-turn."""
    return generated.split(END, 1)[0].strip()
