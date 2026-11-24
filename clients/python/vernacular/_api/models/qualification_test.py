from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from uuid import UUID

if TYPE_CHECKING:
    from ..models.qualification_question import QualificationQuestion


T = TypeVar("T", bound="QualificationTest")


@_attrs_define
class QualificationTest:
    """
    Attributes:
        id (UUID):
        language (str): ISO 639-3 code Example: bem.
        questions (list['QualificationQuestion']):
    """

    id: UUID
    language: str
    questions: list["QualificationQuestion"]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = str(self.id)

        language = self.language

        questions = []
        for questions_item_data in self.questions:
            questions_item = questions_item_data.to_dict()
            questions.append(questions_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "id": id,
                "language": language,
                "questions": questions,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.qualification_question import QualificationQuestion

        d = dict(src_dict)
        id = UUID(d.pop("id"))

        language = d.pop("language")

        questions = []
        _questions = d.pop("questions")
        for questions_item_data in _questions:
            questions_item = QualificationQuestion.from_dict(questions_item_data)

            questions.append(questions_item)

        qualification_test = cls(
            id=id,
            language=language,
            questions=questions,
        )

        qualification_test.additional_properties = d
        return qualification_test

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
