from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from ..models.resolution_outcome import ResolutionOutcome
from ..models.trust_status import TrustStatus
from uuid import UUID


T = TypeVar("T", bound="Resolution")


@_attrs_define
class Resolution:
    """
    Attributes:
        item_id (UUID):
        outcome (ResolutionOutcome):
        status (TrustStatus): How far an item has been checked (ADR-0004). Only `verified` is returned by default;
            `disputed` is never returned.
    """

    item_id: UUID
    outcome: ResolutionOutcome
    status: TrustStatus
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        item_id = str(self.item_id)

        outcome = self.outcome.value

        status = self.status.value

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "item_id": item_id,
                "outcome": outcome,
                "status": status,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        item_id = UUID(d.pop("item_id"))

        outcome = ResolutionOutcome(d.pop("outcome"))

        status = TrustStatus(d.pop("status"))

        resolution = cls(
            item_id=item_id,
            outcome=outcome,
            status=status,
        )

        resolution.additional_properties = d
        return resolution

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
