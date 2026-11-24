from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from ..models.allowed_use import AllowedUse
from ..models.source_kind import SourceKind
from dateutil.parser import isoparse
from typing import Union
from uuid import UUID
import datetime


T = TypeVar("T", bound="Source")


@_attrs_define
class Source:
    """
    Attributes:
        kind (SourceKind): `licensed`: a third party's data used under licence. `contributor`: our volunteers.
            `machine`: an AI provider's output.
        name (str):  Example: GATITOS.
        rights_holder (str):  Example: Google LLC.
        licence (str):  Example: CC BY 4.0.
        allowed_uses (list[AllowedUse]): What the licence allows
        id (UUID):
        has_permission_document (bool):
        created_at (datetime.datetime):
        updated_at (datetime.datetime):
        url (Union[Unset, str]):
        notes (Union[Unset, str]):
    """

    kind: SourceKind
    name: str
    rights_holder: str
    licence: str
    allowed_uses: list[AllowedUse]
    id: UUID
    has_permission_document: bool
    created_at: datetime.datetime
    updated_at: datetime.datetime
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

        id = str(self.id)

        has_permission_document = self.has_permission_document

        created_at = self.created_at.isoformat()

        updated_at = self.updated_at.isoformat()

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
                "id": id,
                "has_permission_document": has_permission_document,
                "created_at": created_at,
                "updated_at": updated_at,
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

        id = UUID(d.pop("id"))

        has_permission_document = d.pop("has_permission_document")

        created_at = isoparse(d.pop("created_at"))

        updated_at = isoparse(d.pop("updated_at"))

        url = d.pop("url", UNSET)

        notes = d.pop("notes", UNSET)

        source = cls(
            kind=kind,
            name=name,
            rights_holder=rights_holder,
            licence=licence,
            allowed_uses=allowed_uses,
            id=id,
            has_permission_document=has_permission_document,
            created_at=created_at,
            updated_at=updated_at,
            url=url,
            notes=notes,
        )

        source.additional_properties = d
        return source

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
