from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field


T = TypeVar("T", bound="ProviderQuality")


@_attrs_define
class ProviderQuality:
    """
    Attributes:
        provider (str): The machine source's name in the register
        language (str): ISO 639-3 code Example: bem.
        decided (int): Review decisions on its candidates
        accepted (int):
        corrected (int):
    """

    provider: str
    language: str
    decided: int
    accepted: int
    corrected: int
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        provider = self.provider

        language = self.language

        decided = self.decided

        accepted = self.accepted

        corrected = self.corrected

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "provider": provider,
                "language": language,
                "decided": decided,
                "accepted": accepted,
                "corrected": corrected,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        provider = d.pop("provider")

        language = d.pop("language")

        decided = d.pop("decided")

        accepted = d.pop("accepted")

        corrected = d.pop("corrected")

        provider_quality = cls(
            provider=provider,
            language=language,
            decided=decided,
            accepted=accepted,
            corrected=corrected,
        )

        provider_quality.additional_properties = d
        return provider_quality

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
