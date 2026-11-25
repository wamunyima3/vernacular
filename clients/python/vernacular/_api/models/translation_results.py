from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from typing import Union

if TYPE_CHECKING:
    from ..models.translation import Translation
    from ..models.machine_translation import MachineTranslation


T = TypeVar("T", bound="TranslationResults")


@_attrs_define
class TranslationResults:
    """
    Attributes:
        data (list['Translation']):
        machine (Union[Unset, MachineTranslation]): Output of a translation model, not checked by any person (TRN-3).
            The model is
            fine-tuned from Gemma: using its output is subject to the Gemma Prohibited Use
            Policy, as set out in the terms of use at /docs/terms.
    """

    data: list["Translation"]
    machine: Union[Unset, "MachineTranslation"] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        data = []
        for data_item_data in self.data:
            data_item = data_item_data.to_dict()
            data.append(data_item)

        machine: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.machine, Unset):
            machine = self.machine.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "data": data,
            }
        )
        if machine is not UNSET:
            field_dict["machine"] = machine

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.translation import Translation
        from ..models.machine_translation import MachineTranslation

        d = dict(src_dict)
        data = []
        _data = d.pop("data")
        for data_item_data in _data:
            data_item = Translation.from_dict(data_item_data)

            data.append(data_item)

        _machine = d.pop("machine", UNSET)
        machine: Union[Unset, MachineTranslation]
        if isinstance(_machine, Unset):
            machine = UNSET
        else:
            machine = MachineTranslation.from_dict(_machine)

        translation_results = cls(
            data=data,
            machine=machine,
        )

        translation_results.additional_properties = d
        return translation_results

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
