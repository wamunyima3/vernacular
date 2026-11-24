from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


if TYPE_CHECKING:
    from ..models.credit import Credit


T = TypeVar("T", bound="CreditList")


@_attrs_define
class CreditList:
    """
    Attributes:
        credits_ (list['Credit']):
    """

    credits_: list["Credit"]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        credits_ = []
        for credits_item_data in self.credits_:
            credits_item = credits_item_data.to_dict()
            credits_.append(credits_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "credits": credits_,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.credit import Credit

        d = dict(src_dict)
        credits_ = []
        _credits_ = d.pop("credits")
        for credits_item_data in _credits_:
            credits_item = Credit.from_dict(credits_item_data)

            credits_.append(credits_item)

        credit_list = cls(
            credits_=credits_,
        )

        credit_list.additional_properties = d
        return credit_list

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
