from enum import Enum


class MinStatus(str, Enum):
    MACHINE = "machine"
    REVIEWED = "reviewed"
    SOURCED = "sourced"
    VERIFIED = "verified"

    def __str__(self) -> str:
        return str(self.value)
