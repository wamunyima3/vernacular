from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define

from ..types import UNSET, Unset

from typing import Union


T = TypeVar("T", bound="RecheckSampleInput")


@_attrs_define
class RecheckSampleInput:
    """
    Attributes:
        quarter (Union[Unset, str]): A calendar quarter, UTC Example: 2026-Q4.
    """

    quarter: Union[Unset, str] = UNSET

    def to_dict(self) -> dict[str, Any]:
        quarter = self.quarter

        field_dict: dict[str, Any] = {}

        field_dict.update({})
        if quarter is not UNSET:
            field_dict["quarter"] = quarter

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        quarter = d.pop("quarter", UNSET)

        recheck_sample_input = cls(
            quarter=quarter,
        )

        return recheck_sample_input
