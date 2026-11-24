from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define

from ..types import UNSET, Unset

from ..models.recheck_outcome import RecheckOutcome
from typing import Union


T = TypeVar("T", bound="RecheckInput")


@_attrs_define
class RecheckInput:
    """
    Attributes:
        outcome (RecheckOutcome):
        note (Union[Unset, str]): What is wrong; needed for `fail`
    """

    outcome: RecheckOutcome
    note: Union[Unset, str] = UNSET

    def to_dict(self) -> dict[str, Any]:
        outcome = self.outcome.value

        note = self.note

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "outcome": outcome,
            }
        )
        if note is not UNSET:
            field_dict["note"] = note

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        outcome = RecheckOutcome(d.pop("outcome"))

        note = d.pop("note", UNSET)

        recheck_input = cls(
            outcome=outcome,
            note=note,
        )

        return recheck_input
