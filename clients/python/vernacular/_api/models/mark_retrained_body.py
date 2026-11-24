from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define


T = TypeVar("T", bound="MarkRetrainedBody")


@_attrs_define
class MarkRetrainedBody:
    """
    Attributes:
        language (str): ISO 639-3 code Example: bem.
    """

    language: str

    def to_dict(self) -> dict[str, Any]:
        language = self.language

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "language": language,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        language = d.pop("language")

        mark_retrained_body = cls(
            language=language,
        )

        return mark_retrained_body
