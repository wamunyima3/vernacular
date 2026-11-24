from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from dateutil.parser import isoparse
from typing import cast
from typing import Union
import datetime


T = TypeVar("T", bound="Speaker")


@_attrs_define
class Speaker:
    """A speaker's consent to be recorded

    Attributes:
        user_id (str):
        languages (list[str]): The languages the speaker may record
        consented_at (datetime.datetime):
        withdrawn_at (Union[Unset, datetime.datetime]): Set once the speaker withdraws; they cannot record until a new
            form is recorded
    """

    user_id: str
    languages: list[str]
    consented_at: datetime.datetime
    withdrawn_at: Union[Unset, datetime.datetime] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        user_id = self.user_id

        languages = self.languages

        consented_at = self.consented_at.isoformat()

        withdrawn_at: Union[Unset, str] = UNSET
        if not isinstance(self.withdrawn_at, Unset):
            withdrawn_at = self.withdrawn_at.isoformat()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "user_id": user_id,
                "languages": languages,
                "consented_at": consented_at,
            }
        )
        if withdrawn_at is not UNSET:
            field_dict["withdrawn_at"] = withdrawn_at

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        user_id = d.pop("user_id")

        languages = cast(list[str], d.pop("languages"))

        consented_at = isoparse(d.pop("consented_at"))

        _withdrawn_at = d.pop("withdrawn_at", UNSET)
        withdrawn_at: Union[Unset, datetime.datetime]
        if isinstance(_withdrawn_at, Unset):
            withdrawn_at = UNSET
        else:
            withdrawn_at = isoparse(_withdrawn_at)

        speaker = cls(
            user_id=user_id,
            languages=languages,
            consented_at=consented_at,
            withdrawn_at=withdrawn_at,
        )

        speaker.additional_properties = d
        return speaker

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
