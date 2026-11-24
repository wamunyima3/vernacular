from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from ..models.review_decision import ReviewDecision
from typing import Union
from uuid import UUID


T = TypeVar("T", bound="Review")


@_attrs_define
class Review:
    """
    Attributes:
        id (UUID):
        item_id (UUID):
        decision (ReviewDecision): `correct`: none of the candidates is right, and the reviewer writes the right word in
            `correction`.
        correction_item_id (Union[Unset, UUID]): The new candidate a correction created
    """

    id: UUID
    item_id: UUID
    decision: ReviewDecision
    correction_item_id: Union[Unset, UUID] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = str(self.id)

        item_id = str(self.item_id)

        decision = self.decision.value

        correction_item_id: Union[Unset, str] = UNSET
        if not isinstance(self.correction_item_id, Unset):
            correction_item_id = str(self.correction_item_id)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "id": id,
                "item_id": item_id,
                "decision": decision,
            }
        )
        if correction_item_id is not UNSET:
            field_dict["correction_item_id"] = correction_item_id

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        id = UUID(d.pop("id"))

        item_id = UUID(d.pop("item_id"))

        decision = ReviewDecision(d.pop("decision"))

        _correction_item_id = d.pop("correction_item_id", UNSET)
        correction_item_id: Union[Unset, UUID]
        if isinstance(_correction_item_id, Unset):
            correction_item_id = UNSET
        else:
            correction_item_id = UUID(_correction_item_id)

        review = cls(
            id=id,
            item_id=item_id,
            decision=decision,
            correction_item_id=correction_item_id,
        )

        review.additional_properties = d
        return review

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
