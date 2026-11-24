from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from ..models.review_card_kind import ReviewCardKind
from ..models.review_priority import ReviewPriority
from uuid import UUID

if TYPE_CHECKING:
    from ..models.review_candidate import ReviewCandidate


T = TypeVar("T", bound="ReviewCard")


@_attrs_define
class ReviewCard:
    """
    Attributes:
        language (str): ISO 639-3 code Example: bem.
        kind (ReviewCardKind): `sense`: candidate words for the meaning. `phrase`: candidate phrases for a line
            of a phrase pack, whose English is `gloss_en`. `example`: sentences using words
            for the meaning, each judged on its own. A card never mixes kinds.
        concept_id (UUID):
        gloss_en (str):  Example: water.
        priority (ReviewPriority): `quick`: every source agrees. `standard`: at least two agree. `full`: no agreement
            (AI-3).
        candidates (list['ReviewCandidate']):
    """

    language: str
    kind: ReviewCardKind
    concept_id: UUID
    gloss_en: str
    priority: ReviewPriority
    candidates: list["ReviewCandidate"]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        language = self.language

        kind = self.kind.value

        concept_id = str(self.concept_id)

        gloss_en = self.gloss_en

        priority = self.priority.value

        candidates = []
        for candidates_item_data in self.candidates:
            candidates_item = candidates_item_data.to_dict()
            candidates.append(candidates_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "language": language,
                "kind": kind,
                "concept_id": concept_id,
                "gloss_en": gloss_en,
                "priority": priority,
                "candidates": candidates,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.review_candidate import ReviewCandidate

        d = dict(src_dict)
        language = d.pop("language")

        kind = ReviewCardKind(d.pop("kind"))

        concept_id = UUID(d.pop("concept_id"))

        gloss_en = d.pop("gloss_en")

        priority = ReviewPriority(d.pop("priority"))

        candidates = []
        _candidates = d.pop("candidates")
        for candidates_item_data in _candidates:
            candidates_item = ReviewCandidate.from_dict(candidates_item_data)

            candidates.append(candidates_item)

        review_card = cls(
            language=language,
            kind=kind,
            concept_id=concept_id,
            gloss_en=gloss_en,
            priority=priority,
            candidates=candidates,
        )

        review_card.additional_properties = d
        return review_card

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
