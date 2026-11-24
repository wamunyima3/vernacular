from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from uuid import UUID


T = TypeVar("T", bound="QualificationResult")


@_attrs_define
class QualificationResult:
    """
    Attributes:
        id (UUID):
        score (int):
        total (int):
        pass_mark (int): Right answers needed to pass
        passed (bool): A senior reviewer of the language confirms a pass before you can review
    """

    id: UUID
    score: int
    total: int
    pass_mark: int
    passed: bool
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = str(self.id)

        score = self.score

        total = self.total

        pass_mark = self.pass_mark

        passed = self.passed

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "id": id,
                "score": score,
                "total": total,
                "pass_mark": pass_mark,
                "passed": passed,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        id = UUID(d.pop("id"))

        score = d.pop("score")

        total = d.pop("total")

        pass_mark = d.pop("pass_mark")

        passed = d.pop("passed")

        qualification_result = cls(
            id=id,
            score=score,
            total=total,
            pass_mark=pass_mark,
            passed=passed,
        )

        qualification_result.additional_properties = d
        return qualification_result

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
