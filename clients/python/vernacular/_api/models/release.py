from collections.abc import Mapping
from typing import Any, TypeVar, TYPE_CHECKING

from attrs import define as _attrs_define
from attrs import field as _attrs_field


from dateutil.parser import isoparse
import datetime

if TYPE_CHECKING:
    from ..models.release_languages import ReleaseLanguages
    from ..models.release_changes import ReleaseChanges


T = TypeVar("T", bound="Release")


@_attrs_define
class Release:
    """
    Attributes:
        name (str):  Example: 2027.1.
        published_at (datetime.datetime):
        notes (str):
        languages (ReleaseLanguages): Counts per language code
        changes (ReleaseChanges): Items compared with the release before; everything is added in the first
    """

    name: str
    published_at: datetime.datetime
    notes: str
    languages: "ReleaseLanguages"
    changes: "ReleaseChanges"
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        name = self.name

        published_at = self.published_at.isoformat()

        notes = self.notes

        languages = self.languages.to_dict()

        changes = self.changes.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "name": name,
                "published_at": published_at,
                "notes": notes,
                "languages": languages,
                "changes": changes,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.release_languages import ReleaseLanguages
        from ..models.release_changes import ReleaseChanges

        d = dict(src_dict)
        name = d.pop("name")

        published_at = isoparse(d.pop("published_at"))

        notes = d.pop("notes")

        languages = ReleaseLanguages.from_dict(d.pop("languages"))

        changes = ReleaseChanges.from_dict(d.pop("changes"))

        release = cls(
            name=name,
            published_at=published_at,
            notes=notes,
            languages=languages,
            changes=changes,
        )

        release.additional_properties = d
        return release

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
