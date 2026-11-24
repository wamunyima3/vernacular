from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field


T = TypeVar("T", bound="RecheckSample")


@_attrs_define
class RecheckSample:
    """
    Attributes:
        quarter (str): A calendar quarter, UTC Example: 2026-Q4.
        sampled (int): Items drawn by this call
    """

    quarter: str
    sampled: int
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        quarter = self.quarter

        sampled = self.sampled

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "quarter": quarter,
                "sampled": sampled,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        quarter = d.pop("quarter")

        sampled = d.pop("sampled")

        recheck_sample = cls(
            quarter=quarter,
            sampled=sampled,
        )

        recheck_sample.additional_properties = d
        return recheck_sample

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
