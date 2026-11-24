from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from dateutil.parser import isoparse
from typing import Union
from uuid import UUID
import datetime


T = TypeVar("T", bound="Key")


@_attrs_define
class Key:
    """
    Attributes:
        id (UUID):
        name (str):
        prefix (str): The key's first characters, to tell keys apart Example: vk_3fA9.
        created_at (datetime.datetime):
        last_used_at (Union[Unset, datetime.datetime]): To the hour
        revoked_at (Union[Unset, datetime.datetime]):
    """

    id: UUID
    name: str
    prefix: str
    created_at: datetime.datetime
    last_used_at: Union[Unset, datetime.datetime] = UNSET
    revoked_at: Union[Unset, datetime.datetime] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = str(self.id)

        name = self.name

        prefix = self.prefix

        created_at = self.created_at.isoformat()

        last_used_at: Union[Unset, str] = UNSET
        if not isinstance(self.last_used_at, Unset):
            last_used_at = self.last_used_at.isoformat()

        revoked_at: Union[Unset, str] = UNSET
        if not isinstance(self.revoked_at, Unset):
            revoked_at = self.revoked_at.isoformat()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "id": id,
                "name": name,
                "prefix": prefix,
                "created_at": created_at,
            }
        )
        if last_used_at is not UNSET:
            field_dict["last_used_at"] = last_used_at
        if revoked_at is not UNSET:
            field_dict["revoked_at"] = revoked_at

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        id = UUID(d.pop("id"))

        name = d.pop("name")

        prefix = d.pop("prefix")

        created_at = isoparse(d.pop("created_at"))

        _last_used_at = d.pop("last_used_at", UNSET)
        last_used_at: Union[Unset, datetime.datetime]
        if isinstance(_last_used_at, Unset):
            last_used_at = UNSET
        else:
            last_used_at = isoparse(_last_used_at)

        _revoked_at = d.pop("revoked_at", UNSET)
        revoked_at: Union[Unset, datetime.datetime]
        if isinstance(_revoked_at, Unset):
            revoked_at = UNSET
        else:
            revoked_at = isoparse(_revoked_at)

        key = cls(
            id=id,
            name=name,
            prefix=prefix,
            created_at=created_at,
            last_used_at=last_used_at,
            revoked_at=revoked_at,
        )

        key.additional_properties = d
        return key

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
