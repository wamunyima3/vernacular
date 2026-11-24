from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


if TYPE_CHECKING:
    from ..models.missing_example import MissingExample


T = TypeVar("T", bound="ExampleCoverage")


@_attrs_define
class ExampleCoverage:
    """How many core meanings in a language have a verified example (LEX-5)

    Attributes:
        language (str): ISO 639-3 code Example: bem.
        core_senses (int): Verified meanings of core-list concepts
        with_example (int): Of those, the meanings with at least one verified example
        missing (list['MissingExample']): Meanings still without a verified example, by English gloss, at most 100
    """

    language: str
    core_senses: int
    with_example: int
    missing: list["MissingExample"]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        language = self.language

        core_senses = self.core_senses

        with_example = self.with_example

        missing = []
        for missing_item_data in self.missing:
            missing_item = missing_item_data.to_dict()
            missing.append(missing_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "language": language,
                "core_senses": core_senses,
                "with_example": with_example,
                "missing": missing,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.missing_example import MissingExample

        d = dict(src_dict)
        language = d.pop("language")

        core_senses = d.pop("core_senses")

        with_example = d.pop("with_example")

        missing = []
        _missing = d.pop("missing")
        for missing_item_data in _missing:
            missing_item = MissingExample.from_dict(missing_item_data)

            missing.append(missing_item)

        example_coverage = cls(
            language=language,
            core_senses=core_senses,
            with_example=with_example,
            missing=missing,
        )

        example_coverage.additional_properties = d
        return example_coverage

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
