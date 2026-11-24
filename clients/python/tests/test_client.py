import json

import httpx
import pytest

from vernacular import Vernacular, VernacularError


def client(status: int, body: object, seen: list[httpx.Request]) -> Vernacular:
    def handler(request: httpx.Request) -> httpx.Response:
        seen.append(request)
        return httpx.Response(status, json=body)

    return Vernacular(api_key="vk_test", base_url="https://api.test", transport=httpx.MockTransport(handler))


def test_search_sends_the_key_and_query() -> None:
    seen: list[httpx.Request] = []
    with client(200, {"data": []}, seen) as v:
        assert v.search("bem", "amenshi", min_status="machine").data == []
    assert seen[0].headers["X-API-Key"] == "vk_test"
    assert seen[0].url.path == "/v1/entries/search"
    assert dict(seen[0].url.params) == {
        "language": "bem",
        "q": "amenshi",
        "limit": "20",
        "min_status": "machine",
    }


def test_errors_raise_with_the_api_code() -> None:
    seen: list[httpx.Request] = []
    body = {"code": "not_found", "message": "Entry not found", "docs_url": "https://docs.test/errors"}
    with client(404, body, seen) as v, pytest.raises(VernacularError) as e:
        v.entry("0b0e4f4a-1f8e-4c39-8b0a-2a6f4d1c9e11")
    assert (e.value.status, e.value.code, e.value.message) == (404, "not_found", "Entry not found")


def test_translate_posts_from_and_to() -> None:
    seen: list[httpx.Request] = []
    with client(200, {"data": []}, seen) as v:
        v.translate("Good morning", "en", "bem")
    assert json.loads(seen[0].content) == {"text": "Good morning", "from": "en", "to": "bem"}
