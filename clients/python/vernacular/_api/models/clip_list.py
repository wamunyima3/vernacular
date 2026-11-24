from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


if TYPE_CHECKING:
    from ..models.clip import Clip


T = TypeVar("T", bound="ClipList")


@_attrs_define
class ClipList:
    """
    Attributes:
        can_record (bool): The signed-in caller has consent in force that covers the entry's language
        data (list['Clip']):
    """

    can_record: bool
    data: list["Clip"]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        can_record = self.can_record

        data = []
        for data_item_data in self.data:
            data_item = data_item_data.to_dict()
            data.append(data_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "can_record": can_record,
                "data": data,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.clip import Clip

        d = dict(src_dict)
        can_record = d.pop("can_record")

        data = []
        _data = d.pop("data")
        for data_item_data in _data:
            data_item = Clip.from_dict(data_item_data)

            data.append(data_item)

        clip_list = cls(
            can_record=can_record,
            data=data,
        )

        clip_list.additional_properties = d
        return clip_list

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
