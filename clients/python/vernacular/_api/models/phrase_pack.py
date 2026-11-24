from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


if TYPE_CHECKING:
    from ..models.phrase import Phrase


T = TypeVar("T", bound="PhrasePack")


@_attrs_define
class PhrasePack:
    """
    Attributes:
        pack (str):  Example: greetings.
        language (str): ISO 639-3 code Example: bem.
        data (list['Phrase']):
    """

    pack: str
    language: str
    data: list["Phrase"]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        pack = self.pack

        language = self.language

        data = []
        for data_item_data in self.data:
            data_item = data_item_data.to_dict()
            data.append(data_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "pack": pack,
                "language": language,
                "data": data,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.phrase import Phrase

        d = dict(src_dict)
        pack = d.pop("pack")

        language = d.pop("language")

        data = []
        _data = d.pop("data")
        for data_item_data in _data:
            data_item = Phrase.from_dict(data_item_data)

            data.append(data_item)

        phrase_pack = cls(
            pack=pack,
            language=language,
            data=data,
        )

        phrase_pack.additional_properties = d
        return phrase_pack

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
