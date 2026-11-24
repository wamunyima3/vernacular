from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field


T = TypeVar("T", bound="KindCounts")


@_attrs_define
class KindCounts:
    """Verified items in a release, by kind

    Attributes:
        entries (int):
        senses (int):
        examples (int):
        phrases (int):
    """

    entries: int
    senses: int
    examples: int
    phrases: int
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        entries = self.entries

        senses = self.senses

        examples = self.examples

        phrases = self.phrases

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "entries": entries,
                "senses": senses,
                "examples": examples,
                "phrases": phrases,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        entries = d.pop("entries")

        senses = d.pop("senses")

        examples = d.pop("examples")

        phrases = d.pop("phrases")

        kind_counts = cls(
            entries=entries,
            senses=senses,
            examples=examples,
            phrases=phrases,
        )

        kind_counts.additional_properties = d
        return kind_counts

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
