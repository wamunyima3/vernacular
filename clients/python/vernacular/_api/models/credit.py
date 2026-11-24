from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field


T = TypeVar("T", bound="Credit")


@_attrs_define
class Credit:
    """
    Attributes:
        name (str):
        words (int): Their contributed words that are verified
        reviews (int): Their review decisions that matched the settled outcome
    """

    name: str
    words: int
    reviews: int
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        name = self.name

        words = self.words

        reviews = self.reviews

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "name": name,
                "words": words,
                "reviews": reviews,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        name = d.pop("name")

        words = d.pop("words")

        reviews = d.pop("reviews")

        credit = cls(
            name=name,
            words=words,
            reviews=reviews,
        )

        credit.additional_properties = d
        return credit

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
