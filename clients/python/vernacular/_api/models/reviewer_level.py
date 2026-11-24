from enum import Enum


class ReviewerLevel(str, Enum):
    REVIEWER = "reviewer"
    SENIOR = "senior"

    def __str__(self) -> str:
        return str(self.value)
