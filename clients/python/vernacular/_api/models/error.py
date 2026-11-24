from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from ..models.error_code import ErrorCode


T = TypeVar("T", bound="Error")


@_attrs_define
class Error:
    """
    Attributes:
        code (ErrorCode): Stable machine-readable code. Codes are never renamed; new codes may be added,
            so treat unknown values as a generic error.
             Example: not_found.
        message (str):  Example: Entry not found.
        docs_url (str):
    """

    code: ErrorCode
    message: str
    docs_url: str
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        code = self.code.value

        message = self.message

        docs_url = self.docs_url

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "code": code,
                "message": message,
                "docs_url": docs_url,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        code = ErrorCode(d.pop("code"))

        message = d.pop("message")

        docs_url = d.pop("docs_url")

        error = cls(
            code=code,
            message=message,
            docs_url=docs_url,
        )

        error.additional_properties = d
        return error

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
