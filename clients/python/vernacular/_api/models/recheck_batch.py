from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from dateutil.parser import isoparse
from typing import Union
from uuid import UUID
import datetime


T = TypeVar("T", bound="RecheckBatch")


@_attrs_define
class RecheckBatch:
    """
    Attributes:
        language (str): ISO 639-3 code Example: bem.
        source_id (UUID):
        source (str): The source's name
        model (str):
        sampled (int):
        checked (int):
        failed (int):
        fails (bool): More than 2% of the checks failed
        sent_back_at (Union[Unset, datetime.datetime]):
        sent_back_items (Union[Unset, int]): Items sent back to review
    """

    language: str
    source_id: UUID
    source: str
    model: str
    sampled: int
    checked: int
    failed: int
    fails: bool
    sent_back_at: Union[Unset, datetime.datetime] = UNSET
    sent_back_items: Union[Unset, int] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        language = self.language

        source_id = str(self.source_id)

        source = self.source

        model = self.model

        sampled = self.sampled

        checked = self.checked

        failed = self.failed

        fails = self.fails

        sent_back_at: Union[Unset, str] = UNSET
        if not isinstance(self.sent_back_at, Unset):
            sent_back_at = self.sent_back_at.isoformat()

        sent_back_items = self.sent_back_items

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "language": language,
                "source_id": source_id,
                "source": source,
                "model": model,
                "sampled": sampled,
                "checked": checked,
                "failed": failed,
                "fails": fails,
            }
        )
        if sent_back_at is not UNSET:
            field_dict["sent_back_at"] = sent_back_at
        if sent_back_items is not UNSET:
            field_dict["sent_back_items"] = sent_back_items

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        language = d.pop("language")

        source_id = UUID(d.pop("source_id"))

        source = d.pop("source")

        model = d.pop("model")

        sampled = d.pop("sampled")

        checked = d.pop("checked")

        failed = d.pop("failed")

        fails = d.pop("fails")

        _sent_back_at = d.pop("sent_back_at", UNSET)
        sent_back_at: Union[Unset, datetime.datetime]
        if isinstance(_sent_back_at, Unset):
            sent_back_at = UNSET
        else:
            sent_back_at = isoparse(_sent_back_at)

        sent_back_items = d.pop("sent_back_items", UNSET)

        recheck_batch = cls(
            language=language,
            source_id=source_id,
            source=source,
            model=model,
            sampled=sampled,
            checked=checked,
            failed=failed,
            fails=fails,
            sent_back_at=sent_back_at,
            sent_back_items=sent_back_items,
        )

        recheck_batch.additional_properties = d
        return recheck_batch

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
