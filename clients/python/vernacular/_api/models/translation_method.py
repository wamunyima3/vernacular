from enum import Enum


class TranslationMethod(str, Enum):
    VERIFIED_PHRASE = "verified_phrase"

    def __str__(self) -> str:
        return str(self.value)
