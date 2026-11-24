from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from typing import Union
from uuid import UUID

if TYPE_CHECKING:
    from ..models.example import Example
    from ..models.trust import Trust


T = TypeVar("T", bound="Sense")


@_attrs_define
class Sense:
    """
    Attributes:
        id (UUID):
        gloss (str):  Example: water.
        examples (list['Example']):
        trust (Trust): Provenance carried by every item (PRV-1). Item schemas include it as a `trust`
            property; never copy these fields into an item schema.
        definition (Union[Unset, str]):
        concept_id (Union[Unset, UUID]): The language-independent meaning this sense expresses (LEX-7)
    """

    id: UUID
    gloss: str
    examples: list["Example"]
    trust: "Trust"
    definition: Union[Unset, str] = UNSET
    concept_id: Union[Unset, UUID] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = str(self.id)

        gloss = self.gloss

        examples = []
        for examples_item_data in self.examples:
            examples_item = examples_item_data.to_dict()
            examples.append(examples_item)

        trust = self.trust.to_dict()

        definition = self.definition

        concept_id: Union[Unset, str] = UNSET
        if not isinstance(self.concept_id, Unset):
            concept_id = str(self.concept_id)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "id": id,
                "gloss": gloss,
                "examples": examples,
                "trust": trust,
            }
        )
        if definition is not UNSET:
            field_dict["definition"] = definition
        if concept_id is not UNSET:
            field_dict["concept_id"] = concept_id

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.example import Example
        from ..models.trust import Trust

        d = dict(src_dict)
        id = UUID(d.pop("id"))

        gloss = d.pop("gloss")

        examples = []
        _examples = d.pop("examples")
        for examples_item_data in _examples:
            examples_item = Example.from_dict(examples_item_data)

            examples.append(examples_item)

        trust = Trust.from_dict(d.pop("trust"))

        definition = d.pop("definition", UNSET)

        _concept_id = d.pop("concept_id", UNSET)
        concept_id: Union[Unset, UUID]
        if isinstance(_concept_id, Unset):
            concept_id = UNSET
        else:
            concept_id = UUID(_concept_id)

        sense = cls(
            id=id,
            gloss=gloss,
            examples=examples,
            trust=trust,
            definition=definition,
            concept_id=concept_id,
        )

        sense.additional_properties = d
        return sense

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
