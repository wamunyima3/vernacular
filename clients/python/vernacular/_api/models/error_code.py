from enum import Enum


class ErrorCode(str, Enum):
    BAD_REQUEST = "bad_request"
    CONFLICT = "conflict"
    FORBIDDEN = "forbidden"
    INTERNAL = "internal"
    METHOD_NOT_ALLOWED = "method_not_allowed"
    NOT_FOUND = "not_found"
    PAYLOAD_TOO_LARGE = "payload_too_large"
    RATE_LIMITED = "rate_limited"
    UNAUTHORIZED = "unauthorized"
    UNAVAILABLE = "unavailable"

    def __str__(self) -> str:
        return str(self.value)
