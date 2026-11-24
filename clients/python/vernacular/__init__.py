"""Python client for the Vernacular API: human-verified words and phrases in
Zambian languages, each with its trust status and sources.

>>> from vernacular import Vernacular
>>> v = Vernacular(api_key="...", base_url="https://api.example.org")
>>> v.search("bem", "amenshi").data

The named methods cover the public reads, return the generated models in
``vernacular.models`` and raise :class:`VernacularError` when a call fails.
``client`` is the generated client for every other endpoint, in
``vernacular._api.api``.
"""

from __future__ import annotations

import json
from collections.abc import Callable
from typing import Any, TypeVar
from uuid import UUID

import httpx
from typing_extensions import Self

from ._api import models
from ._api.api.detect import detect_language
from ._api.api.languages import list_languages
from ._api.api.lexicon import get_concept, get_entry, lookup, search_entries
from ._api.api.phrases import get_phrase_pack, translate
from ._api.api.releases import list_releases
from ._api.client import AuthenticatedClient
from ._api.types import Response

__all__ = ["Vernacular", "VernacularError", "models"]

T = TypeVar("T")


class VernacularError(Exception):
    """What the API said when a call failed: its error code, message and documentation link."""

    def __init__(self, status: int, code: str, message: str, docs_url: str | None = None) -> None:
        super().__init__(message)
        self.status = status
        self.code = code
        self.message = message
        self.docs_url = docs_url


def _options(**kw: Any) -> dict[str, Any]:
    """Keyword arguments for a generated call: None means not sent, and min_status may be a string."""
    out = {k: v for k, v in kw.items() if v is not None}
    if isinstance(out.get("min_status"), str):
        out["min_status"] = models.MinStatus(out["min_status"])
    return out


def _unwrap(response: Response[Any]) -> Any:
    if 200 <= response.status_code < 300:
        return response.parsed
    try:
        body = json.loads(response.content)
    except ValueError:
        body = {}
    raise VernacularError(
        response.status_code,
        body.get("code", "unknown"),
        body.get("message", response.content.decode(errors="replace")),
        body.get("docs_url"),
    )


def _uuid(id: str | UUID) -> UUID:
    try:
        return id if isinstance(id, UUID) else UUID(id)
    except ValueError:
        raise VernacularError(400, "invalid_request", f"{id!r} is not an ID") from None


class Vernacular:
    """A client for the Vernacular API. Use it as a context manager, or call close()."""

    def __init__(self, api_key: str, base_url: str, *, timeout: float = 10.0, **httpx_args: Any) -> None:
        self.client = AuthenticatedClient(
            base_url=base_url,
            token=api_key,
            prefix="",
            auth_header_name="X-API-Key",
            timeout=httpx.Timeout(timeout),
            raise_on_unexpected_status=False,
            httpx_args=httpx_args,
        )

    def __enter__(self) -> Self:
        return self

    def __exit__(self, *_: object) -> None:
        self.close()

    def close(self) -> None:
        self.client.get_httpx_client().close()

    def _call(self, op: Callable[..., Response[T]], **kw: Any) -> T:
        return _unwrap(op(client=self.client, **_options(**kw)))

    def languages(self) -> models.LanguageList:
        """The languages served, with their status (GET /v1/languages)."""
        return self._call(list_languages.sync_detailed)

    def search(self, language: str, q: str, **options: Any) -> models.EntrySearchResults:
        """Search entries by word, variant or stem (GET /v1/entries/search). Options:
        limit, min_status, variety, release."""
        return self._call(search_entries.sync_detailed, language=language, q=q, **options)

    def entry(self, id: str | UUID, **options: Any) -> models.Entry:
        """One entry with its senses, examples and trust fields (GET /v1/entries/{id})."""
        return self._call(get_entry.sync_detailed, id=_uuid(id), **options)

    def concept(self, id: str | UUID, **options: Any) -> models.Concept:
        """A concept and its words across languages (GET /v1/concepts/{id})."""
        return self._call(get_concept.sync_detailed, id=_uuid(id), **options)

    def lookup(self, q: str, from_: str, to: list[str] | None = None, **options: Any) -> models.LookupResults:
        """Words for English glosses, in one or more languages (GET /v1/lookup)."""
        return self._call(lookup.sync_detailed, q=q, from_=from_, to=to, **options)

    def phrases(self, pack: str, lang: str, **options: Any) -> models.PhrasePack:
        """A phrase pack in one language (GET /v1/phrases)."""
        return self._call(get_phrase_pack.sync_detailed, pack=pack, lang=lang, **options)

    def translate(
        self, text: str, from_: str, to: str, *, machine: bool = False, **options: Any
    ) -> models.TranslationResults:
        """Verified phrase translations (POST /v1/translate). With machine=True and no verified
        match, ``machine`` may hold a machine translation, marked as such with a confidence."""
        body = models.TranslateBody.from_dict({"text": text, "from": from_, "to": to, "machine": machine})
        return self._call(translate.sync_detailed, body=body, **options)

    def detect(self, text: str) -> models.Detection:
        """Which language a text is in (POST /v1/detect)."""
        return self._call(
            detect_language.sync_detailed, body=models.DetectLanguageBody.from_dict({"text": text})
        )

    def releases(self) -> models.ReleaseList:
        """Published dataset releases (GET /v1/releases)."""
        return self._call(list_releases.sync_detailed)
