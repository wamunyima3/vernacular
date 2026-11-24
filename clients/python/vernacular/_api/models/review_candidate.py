from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from ..models.part_of_speech import PartOfSpeech
from typing import Union
from uuid import UUID

if TYPE_CHECKING:
    from ..models.corrected_meaning import CorrectedMeaning


T = TypeVar("T", bound="ReviewCandidate")


@_attrs_define
class ReviewCandidate:
    """
    Attributes:
        item_id (UUID): The candidate's item, a sense or, on a phrase card, a phrase
        headword (str): The word, or on a phrase card the phrase Example: amenshi.
        notes (Union[Unset, str]): On a phrase card, when to use this form
        part_of_speech (Union[Unset, PartOfSpeech]):
        noun_class (Union[Unset, str]):  Example: 6.
        plural (Union[Unset, str]):
        corrects (Union[Unset, CorrectedMeaning]): The meaning a correction would replace; verifying the correction
            retires it
    """

    item_id: UUID
    headword: str
    notes: Union[Unset, str] = UNSET
    part_of_speech: Union[Unset, PartOfSpeech] = UNSET
    noun_class: Union[Unset, str] = UNSET
    plural: Union[Unset, str] = UNSET
    corrects: Union[Unset, "CorrectedMeaning"] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        item_id = str(self.item_id)

        headword = self.headword

        notes = self.notes

        part_of_speech: Union[Unset, str] = UNSET
        if not isinstance(self.part_of_speech, Unset):
            part_of_speech = self.part_of_speech.value

        noun_class = self.noun_class

        plural = self.plural

        corrects: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.corrects, Unset):
            corrects = self.corrects.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "item_id": item_id,
                "headword": headword,
            }
        )
        if notes is not UNSET:
            field_dict["notes"] = notes
        if part_of_speech is not UNSET:
            field_dict["part_of_speech"] = part_of_speech
        if noun_class is not UNSET:
            field_dict["noun_class"] = noun_class
        if plural is not UNSET:
            field_dict["plural"] = plural
        if corrects is not UNSET:
            field_dict["corrects"] = corrects

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.corrected_meaning import CorrectedMeaning

        d = dict(src_dict)
        item_id = UUID(d.pop("item_id"))

        headword = d.pop("headword")

        notes = d.pop("notes", UNSET)

        _part_of_speech = d.pop("part_of_speech", UNSET)
        part_of_speech: Union[Unset, PartOfSpeech]
        if isinstance(_part_of_speech, Unset):
            part_of_speech = UNSET
        else:
            part_of_speech = PartOfSpeech(_part_of_speech)

        noun_class = d.pop("noun_class", UNSET)

        plural = d.pop("plural", UNSET)

        _corrects = d.pop("corrects", UNSET)
        corrects: Union[Unset, CorrectedMeaning]
        if isinstance(_corrects, Unset):
            corrects = UNSET
        else:
            corrects = CorrectedMeaning.from_dict(_corrects)

        review_candidate = cls(
            item_id=item_id,
            headword=headword,
            notes=notes,
            part_of_speech=part_of_speech,
            noun_class=noun_class,
            plural=plural,
            corrects=corrects,
        )

        review_candidate.additional_properties = d
        return review_candidate

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
