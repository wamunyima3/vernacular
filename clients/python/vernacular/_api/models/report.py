from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from ..models.trust_status import TrustStatus
from dateutil.parser import isoparse
from uuid import UUID
import datetime


T = TypeVar("T", bound="Report")


@_attrs_define
class Report:
    """
    Attributes:
        id (UUID):
        sense_id (UUID):
        language (str): ISO 639-3 code Example: bem.
        headword (str):
        gloss (str):
        status (TrustStatus): How far an item has been checked (ADR-0004). Only `verified` is returned by default;
            `disputed` is never returned.
        reason (str):
        created_at (datetime.datetime):
    """

    id: UUID
    sense_id: UUID
    language: str
    headword: str
    gloss: str
    status: TrustStatus
    reason: str
    created_at: datetime.datetime
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = str(self.id)

        sense_id = str(self.sense_id)

        language = self.language

        headword = self.headword

        gloss = self.gloss

        status = self.status.value

        reason = self.reason

        created_at = self.created_at.isoformat()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "id": id,
                "sense_id": sense_id,
                "language": language,
                "headword": headword,
                "gloss": gloss,
                "status": status,
                "reason": reason,
                "created_at": created_at,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        id = UUID(d.pop("id"))

        sense_id = UUID(d.pop("sense_id"))

        language = d.pop("language")

        headword = d.pop("headword")

        gloss = d.pop("gloss")

        status = TrustStatus(d.pop("status"))

        reason = d.pop("reason")

        created_at = isoparse(d.pop("created_at"))

        report = cls(
            id=id,
            sense_id=sense_id,
            language=language,
            headword=headword,
            gloss=gloss,
            status=status,
            reason=reason,
            created_at=created_at,
        )

        report.additional_properties = d
        return report

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
