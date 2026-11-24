from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define


if TYPE_CHECKING:
    from ..models.example_input import ExampleInput


T = TypeVar("T", bound="ExampleContributionInput")


@_attrs_define
class ExampleContributionInput:
    """
    Attributes:
        example (ExampleInput): A sentence using the word, with its English translation (LEX-5)
        agree_to_licence (bool): The contributor agrees their example is published under the dataset licence
    """

    example: "ExampleInput"
    agree_to_licence: bool

    def to_dict(self) -> dict[str, Any]:
        example = self.example.to_dict()

        agree_to_licence = self.agree_to_licence

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "example": example,
                "agree_to_licence": agree_to_licence,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.example_input import ExampleInput

        d = dict(src_dict)
        example = ExampleInput.from_dict(d.pop("example"))

        agree_to_licence = d.pop("agree_to_licence")

        example_contribution_input = cls(
            example=example,
            agree_to_licence=agree_to_licence,
        )

        return example_contribution_input
