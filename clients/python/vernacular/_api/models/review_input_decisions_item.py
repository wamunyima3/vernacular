from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define


from ..models.review_decision import ReviewDecision
from uuid import UUID


T = TypeVar("T", bound="ReviewInputDecisionsItem")


@_attrs_define
class ReviewInputDecisionsItem:
    """
    Attributes:
        item_id (UUID):
        decision (ReviewDecision): `correct`: none of the candidates is right, and the reviewer writes the right word in
            `correction`.
    """

    item_id: UUID
    decision: ReviewDecision

    def to_dict(self) -> dict[str, Any]:
        item_id = str(self.item_id)

        decision = self.decision.value

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "item_id": item_id,
                "decision": decision,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        item_id = UUID(d.pop("item_id"))

        decision = ReviewDecision(d.pop("decision"))

        review_input_decisions_item = cls(
            item_id=item_id,
            decision=decision,
        )

        return review_input_decisions_item
