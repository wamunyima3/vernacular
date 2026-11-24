from enum import Enum


class AllowedUse(str, Enum):
    COMMERCIAL = "commercial"
    DERIVE = "derive"
    REDISTRIBUTE = "redistribute"
    TRAIN = "train"

    def __str__(self) -> str:
        return str(self.value)
