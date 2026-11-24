from enum import Enum


class LanguageStatus(str, Enum):
    ACTIVE = "active"
    PLANNED = "planned"

    def __str__(self) -> str:
        return str(self.value)
