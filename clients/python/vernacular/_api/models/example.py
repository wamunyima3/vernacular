from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from uuid import UUID

if TYPE_CHECKING:
    from ..models.trust import Trust


T = TypeVar("T", bound="Example")


@_attrs_define
class Example:
    """
    Attributes:
        id (UUID):
        text (str):
        translation_en (str):
        trust (Trust): Provenance carried by every item (PRV-1). Item schemas include it as a `trust`
            property; never copy these fields into an item schema.
    """

    id: UUID
    text: str
    translation_en: str
    trust: "Trust"
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = str(self.id)

        text = self.text

        translation_en = self.translation_en

        trust = self.trust.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "id": id,
                "text": text,
                "translation_en": translation_en,
                "trust": trust,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.trust import Trust

        d = dict(src_dict)
        id = UUID(d.pop("id"))

        text = d.pop("text")

        translation_en = d.pop("translation_en")

        trust = Trust.from_dict(d.pop("trust"))

        example = cls(
            id=id,
            text=text,
            translation_en=translation_en,
            trust=trust,
        )

        example.additional_properties = d
        return example

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
