from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define


from uuid import UUID


T = TypeVar("T", bound="CreateReportBody")


@_attrs_define
class CreateReportBody:
    """
    Attributes:
        sense_id (UUID):
        reason (str):  Example: This word means moon, not water..
    """

    sense_id: UUID
    reason: str

    def to_dict(self) -> dict[str, Any]:
        sense_id = str(self.sense_id)

        reason = self.reason

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "sense_id": sense_id,
                "reason": reason,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        sense_id = UUID(d.pop("sense_id"))

        reason = d.pop("reason")

        create_report_body = cls(
            sense_id=sense_id,
            reason=reason,
        )

        return create_report_body
