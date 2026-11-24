from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define

from ..types import UNSET, Unset

from ..models.part_of_speech import PartOfSpeech
from typing import Union


T = TypeVar("T", bound="ReviewInputCorrection")


@_attrs_define
class ReviewInputCorrection:
    """The reviewer's own word for the card's meaning, or on a phrase card their own
    phrase, required with `correct` decisions. It becomes a new candidate that other
    reviewers check; the candidates it replaces keep their history (US-22).

        Attributes:
            headword (str):
            part_of_speech (Union[Unset, PartOfSpeech]):
            noun_class (Union[Unset, str]):
    """

    headword: str
    part_of_speech: Union[Unset, PartOfSpeech] = UNSET
    noun_class: Union[Unset, str] = UNSET

    def to_dict(self) -> dict[str, Any]:
        headword = self.headword

        part_of_speech: Union[Unset, str] = UNSET
        if not isinstance(self.part_of_speech, Unset):
            part_of_speech = self.part_of_speech.value

        noun_class = self.noun_class

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "headword": headword,
            }
        )
        if part_of_speech is not UNSET:
            field_dict["part_of_speech"] = part_of_speech
        if noun_class is not UNSET:
            field_dict["noun_class"] = noun_class

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        headword = d.pop("headword")

        _part_of_speech = d.pop("part_of_speech", UNSET)
        part_of_speech: Union[Unset, PartOfSpeech]
        if isinstance(_part_of_speech, Unset):
            part_of_speech = UNSET
        else:
            part_of_speech = PartOfSpeech(_part_of_speech)

        noun_class = d.pop("noun_class", UNSET)

        review_input_correction = cls(
            headword=headword,
            part_of_speech=part_of_speech,
            noun_class=noun_class,
        )

        return review_input_correction
