from enum import Enum


class ResolutionOutcome(str, Enum):
    KEEP = "keep"
    SEND_BACK = "send_back"

    def __str__(self) -> str:
        return str(self.value)
