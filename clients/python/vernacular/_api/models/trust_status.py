from enum import Enum


class TrustStatus(str, Enum):
    DISPUTED = "disputed"
    MACHINE = "machine"
    REVIEWED = "reviewed"
    SOURCED = "sourced"
    VERIFIED = "verified"

    def __str__(self) -> str:
        return str(self.value)
