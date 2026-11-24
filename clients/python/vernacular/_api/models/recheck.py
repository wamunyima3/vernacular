from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


if TYPE_CHECKING:
    from ..models.review_candidate import ReviewCandidate


T = TypeVar("T", bound="Recheck")


@_attrs_define
class Recheck:
    """
    Attributes:
        quarter (str): A calendar quarter, UTC Example: 2026-Q4.
        language (str): ISO 639-3 code Example: bem.
        candidate (ReviewCandidate):
    """

    quarter: str
    language: str
    candidate: "ReviewCandidate"
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        quarter = self.quarter

        language = self.language

        candidate = self.candidate.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "quarter": quarter,
                "language": language,
                "candidate": candidate,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.review_candidate import ReviewCandidate

        d = dict(src_dict)
        quarter = d.pop("quarter")

        language = d.pop("language")

        candidate = ReviewCandidate.from_dict(d.pop("candidate"))

        recheck = cls(
            quarter=quarter,
            language=language,
            candidate=candidate,
        )

        recheck.additional_properties = d
        return recheck

    @property
    def additional_keys(self) -> list[str]:
        return list(self.additional_properties.keys())

    def __getitem__(self, key: str) -> Any:
        return self.additional_properties[key]

    def __setitem__(self, key: str, value: Any) -> None:
        self.additional_properties[key] = value

    def __delitem__(self, key: str) -> None:
        del self.additional_properties[key]

    def __contains__(self, key: str) -> bool:
        return key in self.additional_properties
