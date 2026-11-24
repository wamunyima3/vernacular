from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define


T = TypeVar("T", bound="SetCreditBody")


@_attrs_define
class SetCreditBody:
    """
    Attributes:
        credited (bool):
    """

    credited: bool

    def to_dict(self) -> dict[str, Any]:
        credited = self.credited

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "credited": credited,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        credited = d.pop("credited")

        set_credit_body = cls(
            credited=credited,
        )

        return set_credit_body
