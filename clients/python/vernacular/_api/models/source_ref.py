from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

from typing import Union


T = TypeVar("T", bound="SourceRef")


@_attrs_define
class SourceRef:
    """
    Attributes:
        name (str):  Example: GATITOS.
        licence (str):  Example: CC BY 4.0.
        url (Union[Unset, str]):
        reference (Union[Unset, str]): Page, line or row in the source
        model (Union[Unset, str]): Model version, present only for machine-generated items Example: claude-opus-5-5.
    """

    name: str
    licence: str
    url: Union[Unset, str] = UNSET
    reference: Union[Unset, str] = UNSET
    model: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        name = self.name

        licence = self.licence

        url = self.url

        reference = self.reference

        model = self.model

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "name": name,
                "licence": licence,
            }
        )
        if url is not UNSET:
            field_dict["url"] = url
        if reference is not UNSET:
            field_dict["reference"] = reference
        if model is not UNSET:
            field_dict["model"] = model

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        name = d.pop("name")

        licence = d.pop("licence")

        url = d.pop("url", UNSET)

        reference = d.pop("reference", UNSET)

        model = d.pop("model", UNSET)

        source_ref = cls(
            name=name,
            licence=licence,
            url=url,
            reference=reference,
            model=model,
        )

        source_ref.additional_properties = d
        return source_ref

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
