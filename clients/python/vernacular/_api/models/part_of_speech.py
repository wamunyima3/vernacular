from enum import Enum


class PartOfSpeech(str, Enum):
    ADJECTIVE = "adjective"
    ADVERB = "adverb"
    CONJUNCTION = "conjunction"
    IDEOPHONE = "ideophone"
    INTERJECTION = "interjection"
    NOUN = "noun"
    NUMERAL = "numeral"
    PHRASE = "phrase"
    PREPOSITION = "preposition"
    PRONOUN = "pronoun"
    VERB = "verb"

    def __str__(self) -> str:
        return str(self.value)
