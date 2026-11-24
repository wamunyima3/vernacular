from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from ..models.reviewer_level import ReviewerLevel
from dateutil.parser import isoparse
from typing import Union
import datetime


T = TypeVar("T", bound="Reviewer")


@_attrs_define
class Reviewer:
    """
    Attributes:
        user_id (str):
        name (str):
        level (ReviewerLevel):
        decided (int): Decisions on items whose outcome is settled, since the last retraining
        agreed (int): Of those, the decisions that match the outcome
        suspended (bool):
        retrained_at (Union[Unset, datetime.datetime]):
    """

    user_id: str
    name: str
    level: ReviewerLevel
    decided: int
    agreed: int
    suspended: bool
    retrained_at: Union[Unset, datetime.datetime] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        user_id = self.user_id

        name = self.name

        level = self.level.value

        decided = self.decided

        agreed = self.agreed

        suspended = self.suspended

        retrained_at: Union[Unset, str] = UNSET
        if not isinstance(self.retrained_at, Unset):
            retrained_at = self.retrained_at.isoformat()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "user_id": user_id,
                "name": name,
                "level": level,
                "decided": decided,
                "agreed": agreed,
                "suspended": suspended,
            }
        )
        if retrained_at is not UNSET:
            field_dict["retrained_at"] = retrained_at

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        user_id = d.pop("user_id")

        name = d.pop("name")

        level = ReviewerLevel(d.pop("level"))

        decided = d.pop("decided")

        agreed = d.pop("agreed")

        suspended = d.pop("suspended")

        _retrained_at = d.pop("retrained_at", UNSET)
        retrained_at: Union[Unset, datetime.datetime]
        if isinstance(_retrained_at, Unset):
            retrained_at = UNSET
        else:
            retrained_at = isoparse(_retrained_at)

        reviewer = cls(
            user_id=user_id,
            name=name,
            level=level,
            decided=decided,
            agreed=agreed,
            suspended=suspended,
            retrained_at=retrained_at,
        )

        reviewer.additional_properties = d
        return reviewer

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
