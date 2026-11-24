from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define

from ..types import UNSET, Unset

from typing import Union


T = TypeVar("T", bound="TranslateBody")


@_attrs_define
class TranslateBody:
    """
    Attributes:
        text (str):  Example: Good morning.
        from_ (str): ISO 639-3 code Example: bem.
        to (str): ISO 639-3 code Example: bem.
        machine (Union[Unset, bool]): When nothing verified matches, also ask for a machine translation Default: False.
    """

    text: str
    from_: str
    to: str
    machine: Union[Unset, bool] = False

    def to_dict(self) -> dict[str, Any]:
        text = self.text

        from_ = self.from_

        to = self.to

        machine = self.machine

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "text": text,
                "from": from_,
                "to": to,
            }
        )
        if machine is not UNSET:
            field_dict["machine"] = machine

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        text = d.pop("text")

        from_ = d.pop("from")

        to = d.pop("to")

        machine = d.pop("machine", UNSET)

        translate_body = cls(
            text=text,
            from_=from_,
            to=to,
            machine=machine,
        )

        return translate_body
