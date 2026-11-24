from enum import Enum


class ReviewPriority(str, Enum):
    FULL = "full"
    QUICK = "quick"
    STANDARD = "standard"

    def __str__(self) -> str:
        return str(self.value)
