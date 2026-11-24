from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from dateutil.parser import isoparse
from uuid import UUID
import datetime


T = TypeVar("T", bound="PendingQualification")


@_attrs_define
class PendingQualification:
    """
    Attributes:
        id (UUID):
        user_id (str):
        name (str):
        language (str): ISO 639-3 code Example: bem.
        score (int):
        total (int):
        answered_at (datetime.datetime):
    """

    id: UUID
    user_id: str
    name: str
    language: str
    score: int
    total: int
    answered_at: datetime.datetime
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = str(self.id)

        user_id = self.user_id

        name = self.name

        language = self.language

        score = self.score

        total = self.total

        answered_at = self.answered_at.isoformat()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "id": id,
                "user_id": user_id,
                "name": name,
                "language": language,
                "score": score,
                "total": total,
                "answered_at": answered_at,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        id = UUID(d.pop("id"))

        user_id = d.pop("user_id")

        name = d.pop("name")

        language = d.pop("language")

        score = d.pop("score")

        total = d.pop("total")

        answered_at = isoparse(d.pop("answered_at"))

        pending_qualification = cls(
            id=id,
            user_id=user_id,
            name=name,
            language=language,
            score=score,
            total=total,
            answered_at=answered_at,
        )

        pending_qualification.additional_properties = d
        return pending_qualification

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
