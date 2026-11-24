from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from ..models.match import Match

if TYPE_CHECKING:
    from ..models.entry_summary import EntrySummary


T = TypeVar("T", bound="EntryHit")


@_attrs_define
class EntryHit:
    """
    Attributes:
        match (Match): What the query matched (US-05): the headword (or, from English, the gloss), a
            spelling variant, the root of an inflected or prefixed form, or a similar spelling.
        entry (EntrySummary): A dictionary entry with its meanings (LEX-1, LEX-3, LEX-6)
    """

    match: Match
    entry: "EntrySummary"
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        match = self.match.value

        entry = self.entry.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "match": match,
                "entry": entry,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.entry_summary import EntrySummary

        d = dict(src_dict)
        match = Match(d.pop("match"))

        entry = EntrySummary.from_dict(d.pop("entry"))

        entry_hit = cls(
            match=match,
            entry=entry,
        )

        entry_hit.additional_properties = d
        return entry_hit

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
