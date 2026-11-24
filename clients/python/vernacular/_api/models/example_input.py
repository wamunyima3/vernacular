from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define


T = TypeVar("T", bound="ExampleInput")


@_attrs_define
class ExampleInput:
    """A sentence using the word, with its English translation (LEX-5)

    Attributes:
        text (str):
        translation_en (str):
    """

    text: str
    translation_en: str

    def to_dict(self) -> dict[str, Any]:
        text = self.text

        translation_en = self.translation_en

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "text": text,
                "translation_en": translation_en,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        text = d.pop("text")

        translation_en = d.pop("translation_en")

        example_input = cls(
            text=text,
            translation_en=translation_en,
        )

        return example_input
