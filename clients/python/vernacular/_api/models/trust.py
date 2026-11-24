from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from ..models.trust_status import TrustStatus
from typing import Union

if TYPE_CHECKING:
    from ..models.source_ref import SourceRef
    from ..models.review_summary import ReviewSummary


T = TypeVar("T", bound="Trust")


@_attrs_define
class Trust:
    """Provenance carried by every item (PRV-1). Item schemas include it as a `trust`
    property; never copy these fields into an item schema.

        Attributes:
            status (TrustStatus): How far an item has been checked (ADR-0004). Only `verified` is returned by default;
                `disputed` is never returned.
            sources (list['SourceRef']):
            reviews (ReviewSummary): Counts only. Reviewer identities are not exposed per item.
            release (Union[Unset, str]): Dataset release the item belongs to, absent until first released Example: 2027.1.
    """

    status: TrustStatus
    sources: list["SourceRef"]
    reviews: "ReviewSummary"
    release: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        status = self.status.value

        sources = []
        for sources_item_data in self.sources:
            sources_item = sources_item_data.to_dict()
            sources.append(sources_item)

        reviews = self.reviews.to_dict()

        release = self.release

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "status": status,
                "sources": sources,
                "reviews": reviews,
            }
        )
        if release is not UNSET:
            field_dict["release"] = release

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.source_ref import SourceRef
        from ..models.review_summary import ReviewSummary

        d = dict(src_dict)
        status = TrustStatus(d.pop("status"))

        sources = []
        _sources = d.pop("sources")
        for sources_item_data in _sources:
            sources_item = SourceRef.from_dict(sources_item_data)

            sources.append(sources_item)

        reviews = ReviewSummary.from_dict(d.pop("reviews"))

        release = d.pop("release", UNSET)

        trust = cls(
            status=status,
            sources=sources,
            reviews=reviews,
            release=release,
        )

        trust.additional_properties = d
        return trust

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
