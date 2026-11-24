from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from ..models.translation_method import TranslationMethod

if TYPE_CHECKING:
    from ..models.phrase import Phrase


T = TypeVar("T", bound="Translation")


@_attrs_define
class Translation:
    """
    Attributes:
        method (TranslationMethod): How the translation was found. Only verified phrases so far.
        translation (str):  Example: Mwashibukeni.
        phrase (Phrase):
    """

    method: TranslationMethod
    translation: str
    phrase: "Phrase"
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        method = self.method.value

        translation = self.translation

        phrase = self.phrase.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "method": method,
                "translation": translation,
                "phrase": phrase,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.phrase import Phrase

        d = dict(src_dict)
        method = TranslationMethod(d.pop("method"))

        translation = d.pop("translation")

        phrase = Phrase.from_dict(d.pop("phrase"))

        translation = cls(
            method=method,
            translation=translation,
            phrase=phrase,
        )

        translation.additional_properties = d
        return translation

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
