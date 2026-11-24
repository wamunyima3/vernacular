from enum import Enum


class Match(str, Enum):
    HEADWORD = "headword"
    ROOT = "root"
    SIMILAR = "similar"
    VARIANT = "variant"

    def __str__(self) -> str:
        return str(self.value)
