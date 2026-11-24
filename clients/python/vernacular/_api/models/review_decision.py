from enum import Enum


class ReviewDecision(str, Enum):
    ACCEPT = "accept"
    CORRECT = "correct"
    REJECT = "reject"

    def __str__(self) -> str:
        return str(self.value)
