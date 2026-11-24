from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from ..models.review_decision import ReviewDecision
from dateutil.parser import isoparse
from typing import Union
import datetime


T = TypeVar("T", bound="DisputeReview")


@_attrs_define
class DisputeReview:
    """
    Attributes:
        decision (ReviewDecision): `correct`: none of the candidates is right, and the reviewer writes the right word in
            `correction`.
        comment (str):
        created_at (datetime.datetime):
        correction (Union[Unset, str]): The word the reviewer wrote instead, for `correct` decisions
    """

    decision: ReviewDecision
    comment: str
    created_at: datetime.datetime
    correction: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        decision = self.decision.value

        comment = self.comment

        created_at = self.created_at.isoformat()

        correction = self.correction

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "decision": decision,
                "comment": comment,
                "created_at": created_at,
            }
        )
        if correction is not UNSET:
            field_dict["correction"] = correction

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        decision = ReviewDecision(d.pop("decision"))

        comment = d.pop("comment")

        created_at = isoparse(d.pop("created_at"))

        correction = d.pop("correction", UNSET)

        dispute_review = cls(
            decision=decision,
            comment=comment,
            created_at=created_at,
            correction=correction,
        )

        dispute_review.additional_properties = d
        return dispute_review

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
