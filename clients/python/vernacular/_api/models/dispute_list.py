from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


if TYPE_CHECKING:
    from ..models.dispute import Dispute


T = TypeVar("T", bound="DisputeList")


@_attrs_define
class DisputeList:
    """
    Attributes:
        disputes (list['Dispute']):
    """

    disputes: list["Dispute"]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        disputes = []
        for disputes_item_data in self.disputes:
            disputes_item = disputes_item_data.to_dict()
            disputes.append(disputes_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "disputes": disputes,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.dispute import Dispute

        d = dict(src_dict)
        disputes = []
        _disputes = d.pop("disputes")
        for disputes_item_data in _disputes:
            disputes_item = Dispute.from_dict(disputes_item_data)

            disputes.append(disputes_item)

        dispute_list = cls(
            disputes=disputes,
        )

        dispute_list.additional_properties = d
        return dispute_list

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
