from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define


from ..models.reviewer_level import ReviewerLevel


T = TypeVar("T", bound="ApprovalInput")


@_attrs_define
class ApprovalInput:
    """
    Attributes:
        email (str):
        language (str): ISO 639-3 code Example: bem.
        level (ReviewerLevel):
    """

    email: str
    language: str
    level: ReviewerLevel

    def to_dict(self) -> dict[str, Any]:
        email = self.email

        language = self.language

        level = self.level.value

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "email": email,
                "language": language,
                "level": level,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        email = d.pop("email")

        language = d.pop("language")

        level = ReviewerLevel(d.pop("level"))

        approval_input = cls(
            email=email,
            language=language,
            level=level,
        )

        return approval_input
