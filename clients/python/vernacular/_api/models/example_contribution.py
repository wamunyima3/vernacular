from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from uuid import UUID

if TYPE_CHECKING:
    from ..models.trust import Trust


T = TypeVar("T", bound="ExampleContribution")


@_attrs_define
class ExampleContribution:
    """
    Attributes:
        example_id (UUID):
        trust (Trust): Provenance carried by every item (PRV-1). Item schemas include it as a `trust`
            property; never copy these fields into an item schema.
    """

    example_id: UUID
    trust: "Trust"
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        example_id = str(self.example_id)

        trust = self.trust.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "example_id": example_id,
                "trust": trust,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.trust import Trust

        d = dict(src_dict)
        example_id = UUID(d.pop("example_id"))

        trust = Trust.from_dict(d.pop("trust"))

        example_contribution = cls(
            example_id=example_id,
            trust=trust,
        )

        example_contribution.additional_properties = d
        return example_contribution

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
