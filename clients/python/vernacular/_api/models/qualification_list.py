from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


if TYPE_CHECKING:
    from ..models.pending_qualification import PendingQualification


T = TypeVar("T", bound="QualificationList")


@_attrs_define
class QualificationList:
    """
    Attributes:
        qualifications (list['PendingQualification']):
    """

    qualifications: list["PendingQualification"]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        qualifications = []
        for qualifications_item_data in self.qualifications:
            qualifications_item = qualifications_item_data.to_dict()
            qualifications.append(qualifications_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "qualifications": qualifications,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.pending_qualification import PendingQualification

        d = dict(src_dict)
        qualifications = []
        _qualifications = d.pop("qualifications")
        for qualifications_item_data in _qualifications:
            qualifications_item = PendingQualification.from_dict(qualifications_item_data)

            qualifications.append(qualifications_item)

        qualification_list = cls(
            qualifications=qualifications,
        )

        qualification_list.additional_properties = d
        return qualification_list

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
