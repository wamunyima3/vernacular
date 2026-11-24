from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from dateutil.parser import isoparse
import datetime


T = TypeVar("T", bound="KeyUsageDay")


@_attrs_define
class KeyUsageDay:
    """
    Attributes:
        day (datetime.date):
        calls (int):
        errors (int): Calls answered with an error other than 429
        limited (int): Calls refused with 429 for going over the limit
    """

    day: datetime.date
    calls: int
    errors: int
    limited: int
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        day = self.day.isoformat()

        calls = self.calls

        errors = self.errors

        limited = self.limited

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "day": day,
                "calls": calls,
                "errors": errors,
                "limited": limited,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        day = isoparse(d.pop("day")).date()

        calls = d.pop("calls")

        errors = d.pop("errors")

        limited = d.pop("limited")

        key_usage_day = cls(
            day=day,
            calls=calls,
            errors=errors,
            limited=limited,
        )

        key_usage_day.additional_properties = d
        return key_usage_day

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
