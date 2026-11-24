from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from uuid import UUID

if TYPE_CHECKING:
    from ..models.key_usage_day import KeyUsageDay


T = TypeVar("T", bound="KeyUsage")


@_attrs_define
class KeyUsage:
    """
    Attributes:
        key_id (UUID):
        name (str):
        prefix (str):
        limit_per_minute (int): Calls a key may make per minute; 0 means unlimited
        days (list['KeyUsageDay']):
    """

    key_id: UUID
    name: str
    prefix: str
    limit_per_minute: int
    days: list["KeyUsageDay"]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        key_id = str(self.key_id)

        name = self.name

        prefix = self.prefix

        limit_per_minute = self.limit_per_minute

        days = []
        for days_item_data in self.days:
            days_item = days_item_data.to_dict()
            days.append(days_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "key_id": key_id,
                "name": name,
                "prefix": prefix,
                "limit_per_minute": limit_per_minute,
                "days": days,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.key_usage_day import KeyUsageDay

        d = dict(src_dict)
        key_id = UUID(d.pop("key_id"))

        name = d.pop("name")

        prefix = d.pop("prefix")

        limit_per_minute = d.pop("limit_per_minute")

        days = []
        _days = d.pop("days")
        for days_item_data in _days:
            days_item = KeyUsageDay.from_dict(days_item_data)

            days.append(days_item)

        key_usage = cls(
            key_id=key_id,
            name=name,
            prefix=prefix,
            limit_per_minute=limit_per_minute,
            days=days,
        )

        key_usage.additional_properties = d
        return key_usage

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
