from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from ..models.language_status import LanguageStatus

if TYPE_CHECKING:
    from ..models.entry_counts import EntryCounts


T = TypeVar("T", bound="Language")


@_attrs_define
class Language:
    """
    Attributes:
        code (str): ISO 639-3 code Example: bem.
        name (str): English name Example: Bemba.
        own_name (str): Name of the language in the language itself Example: Ichibemba.
        wave (int): Rollout wave from the delivery plan Example: 1.
        status (LanguageStatus):
        entries (EntryCounts): Live counts of entries. An entry counts as verified when one of its senses is, and
            towards the total when it has a sense that can be served (not disputed).
    """

    code: str
    name: str
    own_name: str
    wave: int
    status: LanguageStatus
    entries: "EntryCounts"
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        code = self.code

        name = self.name

        own_name = self.own_name

        wave = self.wave

        status = self.status.value

        entries = self.entries.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "code": code,
                "name": name,
                "own_name": own_name,
                "wave": wave,
                "status": status,
                "entries": entries,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.entry_counts import EntryCounts

        d = dict(src_dict)
        code = d.pop("code")

        name = d.pop("name")

        own_name = d.pop("own_name")

        wave = d.pop("wave")

        status = LanguageStatus(d.pop("status"))

        entries = EntryCounts.from_dict(d.pop("entries"))

        language = cls(
            code=code,
            name=name,
            own_name=own_name,
            wave=wave,
            status=status,
            entries=entries,
        )

        language.additional_properties = d
        return language

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
