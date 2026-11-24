from enum import Enum


class ReviewCardKind(str, Enum):
    EXAMPLE = "example"
    PHRASE = "phrase"
    SENSE = "sense"

    def __str__(self) -> str:
        return str(self.value)
