from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from ..models.machine_translation_status import MachineTranslationStatus


T = TypeVar("T", bound="MachineTranslation")


@_attrs_define
class MachineTranslation:
    """Output of a translation model, not checked by any person (TRN-3). The model is
    fine-tuned from Gemma: using its output is subject to the Gemma Prohibited Use
    Policy, as set out in the terms of use at /docs/terms.

        Attributes:
            translation (str):  Example: Mulale bwino.
            status (MachineTranslationStatus): Always `machine`; never treat it as verified
            confidence (float): The model's own score
            model (str): The model and version that produced it Example: vernacular-mt-0.1.
    """

    translation: str
    status: MachineTranslationStatus
    confidence: float
    model: str
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        translation = self.translation

        status = self.status.value

        confidence = self.confidence

        model = self.model

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "translation": translation,
                "status": status,
                "confidence": confidence,
                "model": model,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        translation = d.pop("translation")

        status = MachineTranslationStatus(d.pop("status"))

        confidence = d.pop("confidence")

        model = d.pop("model")

        machine_translation = cls(
            translation=translation,
            status=status,
            confidence=confidence,
            model=model,
        )

        machine_translation.additional_properties = d
        return machine_translation

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
