from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from uuid import UUID


T = TypeVar("T", bound="CorrectedMeaning")


@_attrs_define
class CorrectedMeaning:
    """The meaning a correction would replace; verifying the correction retires it

    Attributes:
        item_id (UUID):
        headword (str):
        gloss (str):
    """

    item_id: UUID
    headword: str
    gloss: str
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        item_id = str(self.item_id)

        headword = self.headword

        gloss = self.gloss

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "item_id": item_id,
                "headword": headword,
                "gloss": gloss,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        item_id = UUID(d.pop("item_id"))

        headword = d.pop("headword")

        gloss = d.pop("gloss")

        corrected_meaning = cls(
            item_id=item_id,
            headword=headword,
            gloss=gloss,
        )

        corrected_meaning.additional_properties = d
        return corrected_meaning

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
