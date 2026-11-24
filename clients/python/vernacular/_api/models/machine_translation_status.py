from enum import Enum


class MachineTranslationStatus(str, Enum):
    MACHINE = "machine"

    def __str__(self) -> str:
        return str(self.value)
