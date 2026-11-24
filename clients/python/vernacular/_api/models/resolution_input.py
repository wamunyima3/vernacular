from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define


from ..models.resolution_outcome import ResolutionOutcome


T = TypeVar("T", bound="ResolutionInput")


@_attrs_define
class ResolutionInput:
    """
    Attributes:
        outcome (ResolutionOutcome):
        note (str):
    """

    outcome: ResolutionOutcome
    note: str

    def to_dict(self) -> dict[str, Any]:
        outcome = self.outcome.value

        note = self.note

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "outcome": outcome,
                "note": note,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        outcome = ResolutionOutcome(d.pop("outcome"))

        note = d.pop("note")

        resolution_input = cls(
            outcome=outcome,
            note=note,
        )

        return resolution_input
