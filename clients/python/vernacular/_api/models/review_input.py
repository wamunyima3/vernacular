from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define

from ..types import UNSET, Unset

from typing import Union

if TYPE_CHECKING:
    from ..models.review_input_correction import ReviewInputCorrection
    from ..models.review_input_decisions_item import ReviewInputDecisionsItem


T = TypeVar("T", bound="ReviewInput")


@_attrs_define
class ReviewInput:
    """
    Attributes:
        decisions (list['ReviewInputDecisionsItem']):
        comment (Union[Unset, str]): Why, for rejections
        correction (Union[Unset, ReviewInputCorrection]): The reviewer's own word for the card's meaning, or on a phrase
            card their own
            phrase, required with `correct` decisions. It becomes a new candidate that other
            reviewers check; the candidates it replaces keep their history (US-22).
    """

    decisions: list["ReviewInputDecisionsItem"]
    comment: Union[Unset, str] = UNSET
    correction: Union[Unset, "ReviewInputCorrection"] = UNSET

    def to_dict(self) -> dict[str, Any]:
        decisions = []
        for decisions_item_data in self.decisions:
            decisions_item = decisions_item_data.to_dict()
            decisions.append(decisions_item)

        comment = self.comment

        correction: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.correction, Unset):
            correction = self.correction.to_dict()

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "decisions": decisions,
            }
        )
        if comment is not UNSET:
            field_dict["comment"] = comment
        if correction is not UNSET:
            field_dict["correction"] = correction

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.review_input_correction import ReviewInputCorrection
        from ..models.review_input_decisions_item import ReviewInputDecisionsItem

        d = dict(src_dict)
        decisions = []
        _decisions = d.pop("decisions")
        for decisions_item_data in _decisions:
            decisions_item = ReviewInputDecisionsItem.from_dict(decisions_item_data)

            decisions.append(decisions_item)

        comment = d.pop("comment", UNSET)

        _correction = d.pop("correction", UNSET)
        correction: Union[Unset, ReviewInputCorrection]
        if isinstance(_correction, Unset):
            correction = UNSET
        else:
            correction = ReviewInputCorrection.from_dict(_correction)

        review_input = cls(
            decisions=decisions,
            comment=comment,
            correction=correction,
        )

        return review_input
