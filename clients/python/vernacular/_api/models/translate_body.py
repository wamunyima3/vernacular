from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define


T = TypeVar("T", bound="TranslateBody")


@_attrs_define
class TranslateBody:
    """
    Attributes:
        text (str):  Example: Good morning.
        from_ (str): ISO 639-3 code Example: bem.
        to (str): ISO 639-3 code Example: bem.
    """

    text: str
    from_: str
    to: str

    def to_dict(self) -> dict[str, Any]:
        text = self.text

        from_ = self.from_

        to = self.to

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "text": text,
                "from": from_,
                "to": to,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        text = d.pop("text")

        from_ = d.pop("from")

        to = d.pop("to")

        translate_body = cls(
            text=text,
            from_=from_,
            to=to,
        )

        return translate_body
