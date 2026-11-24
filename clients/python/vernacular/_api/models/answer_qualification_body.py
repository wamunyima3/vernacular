from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define


from uuid import UUID


T = TypeVar("T", bound="AnswerQualificationBody")


@_attrs_define
class AnswerQualificationBody:
    """
    Attributes:
        answers (list[UUID]): The chosen option for each question, in order
    """

    answers: list[UUID]

    def to_dict(self) -> dict[str, Any]:
        answers = []
        for answers_item_data in self.answers:
            answers_item = str(answers_item_data)
            answers.append(answers_item)

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "answers": answers,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        answers = []
        _answers = d.pop("answers")
        for answers_item_data in _answers:
            answers_item = UUID(answers_item_data)

            answers.append(answers_item)

        answer_qualification_body = cls(
            answers=answers,
        )

        return answer_qualification_body
