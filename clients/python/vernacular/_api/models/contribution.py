from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from typing import Union
from uuid import UUID

if TYPE_CHECKING:
    from ..models.trust import Trust


T = TypeVar("T", bound="Contribution")


@_attrs_define
class Contribution:
    """
    Attributes:
        entry_id (UUID):
        sense_id (UUID):
        new_sense (bool): False when the word and meaning were already known; the contribution counts as agreement
        trust (Trust): Provenance carried by every item (PRV-1). Item schemas include it as a `trust`
            property; never copy these fields into an item schema.
        example_id (Union[Unset, UUID]):
    """

    entry_id: UUID
    sense_id: UUID
    new_sense: bool
    trust: "Trust"
    example_id: Union[Unset, UUID] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        entry_id = str(self.entry_id)

        sense_id = str(self.sense_id)

        new_sense = self.new_sense

        trust = self.trust.to_dict()

        example_id: Union[Unset, str] = UNSET
        if not isinstance(self.example_id, Unset):
            example_id = str(self.example_id)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "entry_id": entry_id,
                "sense_id": sense_id,
                "new_sense": new_sense,
                "trust": trust,
            }
        )
        if example_id is not UNSET:
            field_dict["example_id"] = example_id

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.trust import Trust

        d = dict(src_dict)
        entry_id = UUID(d.pop("entry_id"))

        sense_id = UUID(d.pop("sense_id"))

        new_sense = d.pop("new_sense")

        trust = Trust.from_dict(d.pop("trust"))

        _example_id = d.pop("example_id", UNSET)
        example_id: Union[Unset, UUID]
        if isinstance(_example_id, Unset):
            example_id = UNSET
        else:
            example_id = UUID(_example_id)

        contribution = cls(
            entry_id=entry_id,
            sense_id=sense_id,
            new_sense=new_sense,
            trust=trust,
            example_id=example_id,
        )

        contribution.additional_properties = d
        return contribution

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
