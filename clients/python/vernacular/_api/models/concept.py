from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from typing import Union
from uuid import UUID

if TYPE_CHECKING:
    from ..models.entry_summary import EntrySummary


T = TypeVar("T", bound="Concept")


@_attrs_define
class Concept:
    """
    Attributes:
        id (UUID):
        gloss_en (str):  Example: water.
        domain (str):  Example: general.
        core (bool): Part of the core vocabulary list (LEX-7)
        entries (list['EntrySummary']): Entries expressing the concept, by language then headword; at most 200
        definition (Union[Unset, str]):
    """

    id: UUID
    gloss_en: str
    domain: str
    core: bool
    entries: list["EntrySummary"]
    definition: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = str(self.id)

        gloss_en = self.gloss_en

        domain = self.domain

        core = self.core

        entries = []
        for entries_item_data in self.entries:
            entries_item = entries_item_data.to_dict()
            entries.append(entries_item)

        definition = self.definition

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "id": id,
                "gloss_en": gloss_en,
                "domain": domain,
                "core": core,
                "entries": entries,
            }
        )
        if definition is not UNSET:
            field_dict["definition"] = definition

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.entry_summary import EntrySummary

        d = dict(src_dict)
        id = UUID(d.pop("id"))

        gloss_en = d.pop("gloss_en")

        domain = d.pop("domain")

        core = d.pop("core")

        entries = []
        _entries = d.pop("entries")
        for entries_item_data in _entries:
            entries_item = EntrySummary.from_dict(entries_item_data)

            entries.append(entries_item)

        definition = d.pop("definition", UNSET)

        concept = cls(
            id=id,
            gloss_en=gloss_en,
            domain=domain,
            core=core,
            entries=entries,
            definition=definition,
        )

        concept.additional_properties = d
        return concept

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
