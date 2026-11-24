from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from ..models.match import Match

if TYPE_CHECKING:
    from ..models.entry_summary import EntrySummary
    from ..models.lookup_hit_concept import LookupHitConcept


T = TypeVar("T", bound="LookupHit")


@_attrs_define
class LookupHit:
    """
    Attributes:
        match (Match): What the query matched (US-05): the headword (or, from English, the gloss), a
            spelling variant, the root of an inflected or prefixed form, or a similar spelling.
        concept (LookupHitConcept):
        entries (list['EntrySummary']): Entries expressing the concept, each with only those senses
    """

    match: Match
    concept: "LookupHitConcept"
    entries: list["EntrySummary"]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        match = self.match.value

        concept = self.concept.to_dict()

        entries = []
        for entries_item_data in self.entries:
            entries_item = entries_item_data.to_dict()
            entries.append(entries_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "match": match,
                "concept": concept,
                "entries": entries,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.entry_summary import EntrySummary
        from ..models.lookup_hit_concept import LookupHitConcept

        d = dict(src_dict)
        match = Match(d.pop("match"))

        concept = LookupHitConcept.from_dict(d.pop("concept"))

        entries = []
        _entries = d.pop("entries")
        for entries_item_data in _entries:
            entries_item = EntrySummary.from_dict(entries_item_data)

            entries.append(entries_item)

        lookup_hit = cls(
            match=match,
            concept=concept,
            entries=entries,
        )

        lookup_hit.additional_properties = d
        return lookup_hit

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
