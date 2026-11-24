from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define


from uuid import UUID


T = TypeVar("T", bound="RecheckBatchInput")


@_attrs_define
class RecheckBatchInput:
    """
    Attributes:
        quarter (str): A calendar quarter, UTC Example: 2026-Q4.
        language (str): ISO 639-3 code Example: bem.
        source_id (UUID):
        model (str): The model version; empty for sources that are not machine
    """

    quarter: str
    language: str
    source_id: UUID
    model: str

    def to_dict(self) -> dict[str, Any]:
        quarter = self.quarter

        language = self.language

        source_id = str(self.source_id)

        model = self.model

        field_dict: dict[str, Any] = {}

        field_dict.update(
            {
                "quarter": quarter,
                "language": language,
                "source_id": source_id,
                "model": model,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        quarter = d.pop("quarter")

        language = d.pop("language")

        source_id = UUID(d.pop("source_id"))

        model = d.pop("model")

        recheck_batch_input = cls(
            quarter=quarter,
            language=language,
            source_id=source_id,
            model=model,
        )

        return recheck_batch_input
