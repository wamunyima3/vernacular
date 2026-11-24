from enum import Enum


class Register(str, Enum):
    ARCHAIC = "archaic"
    FORMAL = "formal"
    INFORMAL = "informal"
    NEUTRAL = "neutral"
    RESPECTFUL = "respectful"
    TOWN = "town"
    VULGAR = "vulgar"

    def __str__(self) -> str:
        return str(self.value)
