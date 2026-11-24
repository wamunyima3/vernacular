from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define

from ..types import UNSET, Unset

from ..models.part_of_speech import PartOfSpeech
from typing import Union
from uuid import UUID

if TYPE_CHECKING:
    from ..models.example_input import ExampleInput


T = TypeVar("T", bound="ContributionInput")


@_attrs_define
class ContributionInput:
    """
    Attributes:
        language (str): ISO 639-3 code of an active language Example: bem.
        headword (str):  Example: amenshi.
        part_of_speech (PartOfSpeech):
        gloss_en (str):  Example: water.
        agree_to_licence (bool): The contributor agrees their contribution is published under the dataset licence
        noun_class (Union[Unset, str]): Singular/plural class pair, for nouns Example: 6.
        example (Union[Unset, ExampleInput]): A sentence using the word, with its English translation (LEX-5)
        corrects (Union[Unset, UUID]): The meaning (sense ID) this contribution corrects, in the same language. The
            correction is reviewed like any contribution; the original stays served and keeps
            its history meanwhile.
    """

    language: str
    headword: str
    part_of_speech: PartOfSpeech
    gloss_en: str
    agree_to_licence: bool
    noun_class: Union[Unset, str] = UNSET
    example: Union[Unset, "ExampleInput"] = UNSET
    corrects: Union[Unset, UUID] = UNSET

    def to_dict(self) -> dict[str, Any]:
        language = self.language

        headword = self.headword

        part_of_speech = self.part_of_speech.value

        gloss_en = self.gloss_en

        agree_to_licence = self.agree_to_licence

        noun_class = self.noun_class

        example: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.example, Unset):
            example = self.example.to_dict()

        corrects: Union[Unset, str] = UNSET
        if not isinstance(self.corrects, Unset):
            corrects = str(self.corrects)

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "language": language,
                "headword": headword,
                "part_of_speech": part_of_speech,
                "gloss_en": gloss_en,
                "agree_to_licence": agree_to_licence,
            }
        )
        if noun_class is not UNSET:
            field_dict["noun_class"] = noun_class
        if example is not UNSET:
            field_dict["example"] = example
        if corrects is not UNSET:
            field_dict["corrects"] = corrects

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.example_input import ExampleInput

        d = dict(src_dict)
        language = d.pop("language")

        headword = d.pop("headword")

        part_of_speech = PartOfSpeech(d.pop("part_of_speech"))

        gloss_en = d.pop("gloss_en")

        agree_to_licence = d.pop("agree_to_licence")

        noun_class = d.pop("noun_class", UNSET)

        _example = d.pop("example", UNSET)
        example: Union[Unset, ExampleInput]
        if isinstance(_example, Unset):
            example = UNSET
        else:
            example = ExampleInput.from_dict(_example)

        _corrects = d.pop("corrects", UNSET)
        corrects: Union[Unset, UUID]
        if isinstance(_corrects, Unset):
            corrects = UNSET
        else:
            corrects = UUID(_corrects)

        contribution_input = cls(
            language=language,
            headword=headword,
            part_of_speech=part_of_speech,
            gloss_en=gloss_en,
            agree_to_licence=agree_to_licence,
            noun_class=noun_class,
            example=example,
            corrects=corrects,
        )

        return contribution_input
