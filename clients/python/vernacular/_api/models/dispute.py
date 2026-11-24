from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from dateutil.parser import isoparse
import datetime

if TYPE_CHECKING:
    from ..models.review_candidate import ReviewCandidate
    from ..models.dispute_review import DisputeReview


T = TypeVar("T", bound="Dispute")


@_attrs_define
class Dispute:
    """
    Attributes:
        language (str): ISO 639-3 code Example: bem.
        gloss_en (str):
        candidate (ReviewCandidate):
        disputed_since (datetime.datetime): When the last counted review was made
        reviews (list['DisputeReview']): Oldest first. Reviewers stay anonymous.
    """

    language: str
    gloss_en: str
    candidate: "ReviewCandidate"
    disputed_since: datetime.datetime
    reviews: list["DisputeReview"]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        language = self.language

        gloss_en = self.gloss_en

        candidate = self.candidate.to_dict()

        disputed_since = self.disputed_since.isoformat()

        reviews = []
        for reviews_item_data in self.reviews:
            reviews_item = reviews_item_data.to_dict()
            reviews.append(reviews_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "language": language,
                "gloss_en": gloss_en,
                "candidate": candidate,
                "disputed_since": disputed_since,
                "reviews": reviews,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.review_candidate import ReviewCandidate
        from ..models.dispute_review import DisputeReview

        d = dict(src_dict)
        language = d.pop("language")

        gloss_en = d.pop("gloss_en")

        candidate = ReviewCandidate.from_dict(d.pop("candidate"))

        disputed_since = isoparse(d.pop("disputed_since"))

        reviews = []
        _reviews = d.pop("reviews")
        for reviews_item_data in _reviews:
            reviews_item = DisputeReview.from_dict(reviews_item_data)

            reviews.append(reviews_item)

        dispute = cls(
            language=language,
            gloss_en=gloss_en,
            candidate=candidate,
            disputed_since=disputed_since,
            reviews=reviews,
        )

        dispute.additional_properties = d
        return dispute

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
