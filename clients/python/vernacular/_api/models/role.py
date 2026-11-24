from enum import Enum


class Role(str, Enum):
    ADMIN = "admin"
    CONTRIBUTOR = "contributor"
    DEVELOPER = "developer"

    def __str__(self) -> str:
        return str(self.value)
