from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from dateutil.parser import isoparse
from typing import Union
import datetime


T = TypeVar("T", bound="ReviewSummary")


@_attrs_define
class ReviewSummary:
    """Counts only. Reviewer identities are not exposed per item.

    Attributes:
        accepts (int):
        corrections (int):
        last_reviewed_at (Union[Unset, datetime.datetime]):
    """

    accepts: int
    corrections: int
    last_reviewed_at: Union[Unset, datetime.datetime] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        accepts = self.accepts

        corrections = self.corrections

        last_reviewed_at: Union[Unset, str] = UNSET
        if not isinstance(self.last_reviewed_at, Unset):
            last_reviewed_at = self.last_reviewed_at.isoformat()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "accepts": accepts,
                "corrections": corrections,
            }
        )
        if last_reviewed_at is not UNSET:
            field_dict["last_reviewed_at"] = last_reviewed_at

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        accepts = d.pop("accepts")

        corrections = d.pop("corrections")

        _last_reviewed_at = d.pop("last_reviewed_at", UNSET)
        last_reviewed_at: Union[Unset, datetime.datetime]
        if isinstance(_last_reviewed_at, Unset):
            last_reviewed_at = UNSET
        else:
            last_reviewed_at = isoparse(_last_reviewed_at)

        review_summary = cls(
            accepts=accepts,
            corrections=corrections,
            last_reviewed_at=last_reviewed_at,
        )

        review_summary.additional_properties = d
        return review_summary

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
