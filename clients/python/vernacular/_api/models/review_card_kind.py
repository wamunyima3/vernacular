from enum import Enum


class ReviewCardKind(str, Enum):
    PHRASE = "phrase"
    SENSE = "sense"

    def __str__(self) -> str:
        return str(self.value)
