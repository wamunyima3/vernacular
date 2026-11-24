from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from typing import cast
from typing import Union
from uuid import UUID

if TYPE_CHECKING:
    from ..models.trust import Trust


T = TypeVar("T", bound="Phrase")


@_attrs_define
class Phrase:
    """
    Attributes:
        id (UUID):
        text (str):  Example: Mwashibukeni.
        source_en (str):  Example: Good morning.
        trust (Trust): Provenance carried by every item (PRV-1). Item schemas include it as a `trust`
            property; never copy these fields into an item schema.
        notes (Union[Unset, str]): Context, and how the words around a placeholder change with what fills it
        placeholders (Union[Unset, list[str]]): Names in braces in `source_en` and `text`, such as `name` for `{name}`,
            which your app
            fills in. Absent when the phrase has none (PHR-2).
        concept_id (Union[Unset, UUID]):
    """

    id: UUID
    text: str
    source_en: str
    trust: "Trust"
    notes: Union[Unset, str] = UNSET
    placeholders: Union[Unset, list[str]] = UNSET
    concept_id: Union[Unset, UUID] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = str(self.id)

        text = self.text

        source_en = self.source_en

        trust = self.trust.to_dict()

        notes = self.notes

        placeholders: Union[Unset, list[str]] = UNSET
        if not isinstance(self.placeholders, Unset):
            placeholders = self.placeholders

        concept_id: Union[Unset, str] = UNSET
        if not isinstance(self.concept_id, Unset):
            concept_id = str(self.concept_id)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "id": id,
                "text": text,
                "source_en": source_en,
                "trust": trust,
            }
        )
        if notes is not UNSET:
            field_dict["notes"] = notes
        if placeholders is not UNSET:
            field_dict["placeholders"] = placeholders
        if concept_id is not UNSET:
            field_dict["concept_id"] = concept_id

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.trust import Trust

        d = dict(src_dict)
        id = UUID(d.pop("id"))

        text = d.pop("text")

        source_en = d.pop("source_en")

        trust = Trust.from_dict(d.pop("trust"))

        notes = d.pop("notes", UNSET)

        placeholders = cast(list[str], d.pop("placeholders", UNSET))

        _concept_id = d.pop("concept_id", UNSET)
        concept_id: Union[Unset, UUID]
        if isinstance(_concept_id, Unset):
            concept_id = UNSET
        else:
            concept_id = UUID(_concept_id)

        phrase = cls(
            id=id,
            text=text,
            source_en=source_en,
            trust=trust,
            notes=notes,
            placeholders=placeholders,
            concept_id=concept_id,
        )

        phrase.additional_properties = d
        return phrase

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
