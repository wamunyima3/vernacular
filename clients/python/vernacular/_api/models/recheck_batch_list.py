from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


if TYPE_CHECKING:
    from ..models.recheck_batch import RecheckBatch


T = TypeVar("T", bound="RecheckBatchList")


@_attrs_define
class RecheckBatchList:
    """
    Attributes:
        quarter (str): A calendar quarter, UTC Example: 2026-Q4.
        batches (list['RecheckBatch']):
    """

    quarter: str
    batches: list["RecheckBatch"]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        quarter = self.quarter

        batches = []
        for batches_item_data in self.batches:
            batches_item = batches_item_data.to_dict()
            batches.append(batches_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "quarter": quarter,
                "batches": batches,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.recheck_batch import RecheckBatch

        d = dict(src_dict)
        quarter = d.pop("quarter")

        batches = []
        _batches = d.pop("batches")
        for batches_item_data in _batches:
            batches_item = RecheckBatch.from_dict(batches_item_data)

            batches.append(batches_item)

        recheck_batch_list = cls(
            quarter=quarter,
            batches=batches,
        )

        recheck_batch_list.additional_properties = d
        return recheck_batch_list

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
