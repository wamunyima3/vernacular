from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from dateutil.parser import isoparse
from uuid import UUID
import datetime


T = TypeVar("T", bound="Clip")


@_attrs_define
class Clip:
    """A pronunciation recorded by a native speaker (LEX-8)

    Attributes:
        id (UUID):
        entry_id (UUID):
        url (str): Where to fetch the sound Example: /v1/audio/0b6c….
        content_type (str):  Example: audio/webm.
        recorded_at (datetime.datetime):
        mine (bool): Recorded by the signed-in caller
    """

    id: UUID
    entry_id: UUID
    url: str
    content_type: str
    recorded_at: datetime.datetime
    mine: bool
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = str(self.id)

        entry_id = str(self.entry_id)

        url = self.url

        content_type = self.content_type

        recorded_at = self.recorded_at.isoformat()

        mine = self.mine

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "id": id,
                "entry_id": entry_id,
                "url": url,
                "content_type": content_type,
                "recorded_at": recorded_at,
                "mine": mine,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        id = UUID(d.pop("id"))

        entry_id = UUID(d.pop("entry_id"))

        url = d.pop("url")

        content_type = d.pop("content_type")

        recorded_at = isoparse(d.pop("recorded_at"))

        mine = d.pop("mine")

        clip = cls(
            id=id,
            entry_id=entry_id,
            url=url,
            content_type=content_type,
            recorded_at=recorded_at,
            mine=mine,
        )

        clip.additional_properties = d
        return clip

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
