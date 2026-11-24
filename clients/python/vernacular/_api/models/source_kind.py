from enum import Enum


class SourceKind(str, Enum):
    CONTRIBUTOR = "contributor"
    LICENSED = "licensed"
    MACHINE = "machine"

    def __str__(self) -> str:
        return str(self.value)
