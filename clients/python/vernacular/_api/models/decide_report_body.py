from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define


T = TypeVar("T", bound="DecideReportBody")


@_attrs_define
class DecideReportBody:
    """
    Attributes:
        confirm (bool):
        note (str):
    """

    confirm: bool
    note: str

    def to_dict(self) -> dict[str, Any]:
        confirm = self.confirm

        note = self.note

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "confirm": confirm,
                "note": note,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        confirm = d.pop("confirm")

        note = d.pop("note")

        decide_report_body = cls(
            confirm=confirm,
            note=note,
        )

        return decide_report_body
