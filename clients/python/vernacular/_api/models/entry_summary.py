from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from ..models.part_of_speech import PartOfSpeech
from ..models.register import Register
from typing import Union
from uuid import UUID

if TYPE_CHECKING:
    from ..models.sense_summary import SenseSummary
    from ..models.trust import Trust


T = TypeVar("T", bound="EntrySummary")


@_attrs_define
class EntrySummary:
    """A dictionary entry with its meanings (LEX-1, LEX-3, LEX-6)

    Attributes:
        id (UUID):
        language (str): ISO 639-3 code Example: bem.
        headword (str):  Example: amenshi.
        register (Register): Social register of a word (LEX-6)
        senses (list['SenseSummary']):
        trust (Trust): Provenance carried by every item (PRV-1). Item schemas include it as a `trust`
            property; never copy these fields into an item schema.
        part_of_speech (Union[Unset, PartOfSpeech]):
        noun_class (Union[Unset, str]):  Example: 6.
        plural (Union[Unset, str]):
        root (Union[Unset, str]):
        variety (Union[Unset, str]): Variety code when not the standard language
    """

    id: UUID
    language: str
    headword: str
    register: Register
    senses: list["SenseSummary"]
    trust: "Trust"
    part_of_speech: Union[Unset, PartOfSpeech] = UNSET
    noun_class: Union[Unset, str] = UNSET
    plural: Union[Unset, str] = UNSET
    root: Union[Unset, str] = UNSET
    variety: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = str(self.id)

        language = self.language

        headword = self.headword

        register = self.register.value

        senses = []
        for senses_item_data in self.senses:
            senses_item = senses_item_data.to_dict()
            senses.append(senses_item)

        trust = self.trust.to_dict()

        part_of_speech: Union[Unset, str] = UNSET
        if not isinstance(self.part_of_speech, Unset):
            part_of_speech = self.part_of_speech.value

        noun_class = self.noun_class

        plural = self.plural

        root = self.root

        variety = self.variety

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "id": id,
                "language": language,
                "headword": headword,
                "register": register,
                "senses": senses,
                "trust": trust,
            }
        )
        if part_of_speech is not UNSET:
            field_dict["part_of_speech"] = part_of_speech
        if noun_class is not UNSET:
            field_dict["noun_class"] = noun_class
        if plural is not UNSET:
            field_dict["plural"] = plural
        if root is not UNSET:
            field_dict["root"] = root
        if variety is not UNSET:
            field_dict["variety"] = variety

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.sense_summary import SenseSummary
        from ..models.trust import Trust

        d = dict(src_dict)
        id = UUID(d.pop("id"))

        language = d.pop("language")

        headword = d.pop("headword")

        register = Register(d.pop("register"))

        senses = []
        _senses = d.pop("senses")
        for senses_item_data in _senses:
            senses_item = SenseSummary.from_dict(senses_item_data)

            senses.append(senses_item)

        trust = Trust.from_dict(d.pop("trust"))

        _part_of_speech = d.pop("part_of_speech", UNSET)
        part_of_speech: Union[Unset, PartOfSpeech]
        if isinstance(_part_of_speech, Unset):
            part_of_speech = UNSET
        else:
            part_of_speech = PartOfSpeech(_part_of_speech)

        noun_class = d.pop("noun_class", UNSET)

        plural = d.pop("plural", UNSET)

        root = d.pop("root", UNSET)

        variety = d.pop("variety", UNSET)

        entry_summary = cls(
            id=id,
            language=language,
            headword=headword,
            register=register,
            senses=senses,
            trust=trust,
            part_of_speech=part_of_speech,
            noun_class=noun_class,
            plural=plural,
            root=root,
            variety=variety,
        )

        entry_summary.additional_properties = d
        return entry_summary

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
