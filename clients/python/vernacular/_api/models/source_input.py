from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from ..models.allowed_use import AllowedUse
from ..models.source_kind import SourceKind
from typing import Union


T = TypeVar("T", bound="SourceInput")


@_attrs_define
class SourceInput:
    """
    Attributes:
        kind (SourceKind): `licensed`: a third party's data used under licence. `contributor`: our volunteers.
            `machine`: an AI provider's output.
        name (str):  Example: GATITOS.
        rights_holder (str):  Example: Google LLC.
        licence (str):  Example: CC BY 4.0.
        allowed_uses (list[AllowedUse]): What the licence allows
        url (Union[Unset, str]):
        notes (Union[Unset, str]):
    """

    kind: SourceKind
    name: str
    rights_holder: str
    licence: str
    allowed_uses: list[AllowedUse]
    url: Union[Unset, str] = UNSET
    notes: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        kind = self.kind.value

        name = self.name

        rights_holder = self.rights_holder

        licence = self.licence

        allowed_uses = []
        for allowed_uses_item_data in self.allowed_uses:
            allowed_uses_item = allowed_uses_item_data.value
            allowed_uses.append(allowed_uses_item)

        url = self.url

        notes = self.notes

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "kind": kind,
                "name": name,
                "rights_holder": rights_holder,
                "licence": licence,
                "allowed_uses": allowed_uses,
            }
        )
        if url is not UNSET:
            field_dict["url"] = url
        if notes is not UNSET:
            field_dict["notes"] = notes

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        kind = SourceKind(d.pop("kind"))

        name = d.pop("name")

        rights_holder = d.pop("rights_holder")

        licence = d.pop("licence")

        allowed_uses = []
        _allowed_uses = d.pop("allowed_uses")
        for allowed_uses_item_data in _allowed_uses:
            allowed_uses_item = AllowedUse(allowed_uses_item_data)

            allowed_uses.append(allowed_uses_item)

        url = d.pop("url", UNSET)

        notes = d.pop("notes", UNSET)

        source_input = cls(
            kind=kind,
            name=name,
            rights_holder=rights_holder,
            licence=licence,
            allowed_uses=allowed_uses,
            url=url,
            notes=notes,
        )

        source_input.additional_properties = d
        return source_input

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
