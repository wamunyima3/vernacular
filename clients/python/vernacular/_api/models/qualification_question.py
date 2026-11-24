from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


if TYPE_CHECKING:
    from ..models.qualification_option import QualificationOption


T = TypeVar("T", bound="QualificationQuestion")


@_attrs_define
class QualificationQuestion:
    """
    Attributes:
        gloss_en (str): The English meaning
        options (list['QualificationOption']): Four words; one means gloss_en
    """

    gloss_en: str
    options: list["QualificationOption"]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        gloss_en = self.gloss_en

        options = []
        for options_item_data in self.options:
            options_item = options_item_data.to_dict()
            options.append(options_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "gloss_en": gloss_en,
                "options": options,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.qualification_option import QualificationOption

        d = dict(src_dict)
        gloss_en = d.pop("gloss_en")

        options = []
        _options = d.pop("options")
        for options_item_data in _options:
            options_item = QualificationOption.from_dict(options_item_data)

            options.append(options_item)

        qualification_question = cls(
            gloss_en=gloss_en,
            options=options,
        )

        qualification_question.additional_properties = d
        return qualification_question

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
