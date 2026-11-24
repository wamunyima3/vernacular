from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from ..models.role import Role

if TYPE_CHECKING:
    from ..models.approval import Approval


T = TypeVar("T", bound="Me")


@_attrs_define
class Me:
    """
    Attributes:
        id (str):
        name (str):
        email (str):
        role (Role): Global account role. Review rights are per language, in `approvals`.
        approvals (list['Approval']):
        credited (bool): The user chose to be credited by name on the public contributors page
    """

    id: str
    name: str
    email: str
    role: Role
    approvals: list["Approval"]
    credited: bool
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = self.id

        name = self.name

        email = self.email

        role = self.role.value

        approvals = []
        for approvals_item_data in self.approvals:
            approvals_item = approvals_item_data.to_dict()
            approvals.append(approvals_item)

        credited = self.credited

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "id": id,
                "name": name,
                "email": email,
                "role": role,
                "approvals": approvals,
                "credited": credited,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.approval import Approval

        d = dict(src_dict)
        id = d.pop("id")

        name = d.pop("name")

        email = d.pop("email")

        role = Role(d.pop("role"))

        approvals = []
        _approvals = d.pop("approvals")
        for approvals_item_data in _approvals:
            approvals_item = Approval.from_dict(approvals_item_data)

            approvals.append(approvals_item)

        credited = d.pop("credited")

        me = cls(
            id=id,
            name=name,
            email=email,
            role=role,
            approvals=approvals,
            credited=credited,
        )

        me.additional_properties = d
        return me

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
