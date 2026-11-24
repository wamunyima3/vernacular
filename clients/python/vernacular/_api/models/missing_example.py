from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from uuid import UUID


T = TypeVar("T", bound="MissingExample")


@_attrs_define
class MissingExample:
    """
    Attributes:
        sense_id (UUID):
        entry_id (UUID):
        headword (str):  Example: amenshi.
        gloss_en (str):  Example: water.
        pending (int): Examples of this meaning waiting for review
    """

    sense_id: UUID
    entry_id: UUID
    headword: str
    gloss_en: str
    pending: int
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        sense_id = str(self.sense_id)

        entry_id = str(self.entry_id)

        headword = self.headword

        gloss_en = self.gloss_en

        pending = self.pending

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "sense_id": sense_id,
                "entry_id": entry_id,
                "headword": headword,
                "gloss_en": gloss_en,
                "pending": pending,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        sense_id = UUID(d.pop("sense_id"))

        entry_id = UUID(d.pop("entry_id"))

        headword = d.pop("headword")

        gloss_en = d.pop("gloss_en")

        pending = d.pop("pending")

        missing_example = cls(
            sense_id=sense_id,
            entry_id=entry_id,
            headword=headword,
            gloss_en=gloss_en,
            pending=pending,
        )

        missing_example.additional_properties = d
        return missing_example

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
