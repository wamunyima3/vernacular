from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.entry_search_results import EntrySearchResults
from ...models.error import Error
from ...models.min_status import MinStatus
from ...types import Unset


def _get_kwargs(
    *,
    language: str,
    q: str,
    limit: Union[Unset, int] = 20,
    min_status: Union[Unset, MinStatus] = UNSET,
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["language"] = language

    params["q"] = q

    params["limit"] = limit

    json_min_status: Union[Unset, str] = UNSET
    if not isinstance(min_status, Unset):
        json_min_status = min_status.value

    params["min_status"] = json_min_status

    params["variety"] = variety

    params["release"] = release

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/v1/entries/search",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[EntrySearchResults, Error]:
    if response.status_code == 200:
        response_200 = EntrySearchResults.from_dict(response.json())

        return response_200

    if response.status_code == 400:
        response_400 = Error.from_dict(response.json())

        return response_400

    if response.status_code == 401:
        response_401 = Error.from_dict(response.json())

        return response_401

    if response.status_code == 429:
        response_429 = Error.from_dict(response.json())

        return response_429

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[EntrySearchResults, Error]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    language: str,
    q: str,
    limit: Union[Unset, int] = 20,
    min_status: Union[Unset, MinStatus] = UNSET,
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Response[Union[EntrySearchResults, Error]]:
    """Search a language's dictionary entries

     Finds entries by headword, tolerating case, accents, spelling variants, missing
    prefixes and small typing errors (LEX-4, NFR-7). Each hit says what matched:
    the headword, a spelling variant, the root (an inflected or prefixed form) or a
    similar spelling. Best matches first. Disputed items are never returned.

    Args:
        language (str): ISO 639-3 code Example: bem.
        q (str):
        limit (Union[Unset, int]):  Default: 20.
        min_status (Union[Unset, MinStatus]):
        variety (Union[Unset, str]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[EntrySearchResults, Error]]
    """

    kwargs = _get_kwargs(
        language=language,
        q=q,
        limit=limit,
        min_status=min_status,
        variety=variety,
        release=release,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: AuthenticatedClient,
    language: str,
    q: str,
    limit: Union[Unset, int] = 20,
    min_status: Union[Unset, MinStatus] = UNSET,
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Optional[Union[EntrySearchResults, Error]]:
    """Search a language's dictionary entries

     Finds entries by headword, tolerating case, accents, spelling variants, missing
    prefixes and small typing errors (LEX-4, NFR-7). Each hit says what matched:
    the headword, a spelling variant, the root (an inflected or prefixed form) or a
    similar spelling. Best matches first. Disputed items are never returned.

    Args:
        language (str): ISO 639-3 code Example: bem.
        q (str):
        limit (Union[Unset, int]):  Default: 20.
        min_status (Union[Unset, MinStatus]):
        variety (Union[Unset, str]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[EntrySearchResults, Error]
    """

    return sync_detailed(
        client=client,
        language=language,
        q=q,
        limit=limit,
        min_status=min_status,
        variety=variety,
        release=release,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    language: str,
    q: str,
    limit: Union[Unset, int] = 20,
    min_status: Union[Unset, MinStatus] = UNSET,
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Response[Union[EntrySearchResults, Error]]:
    """Search a language's dictionary entries

     Finds entries by headword, tolerating case, accents, spelling variants, missing
    prefixes and small typing errors (LEX-4, NFR-7). Each hit says what matched:
    the headword, a spelling variant, the root (an inflected or prefixed form) or a
    similar spelling. Best matches first. Disputed items are never returned.

    Args:
        language (str): ISO 639-3 code Example: bem.
        q (str):
        limit (Union[Unset, int]):  Default: 20.
        min_status (Union[Unset, MinStatus]):
        variety (Union[Unset, str]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[EntrySearchResults, Error]]
    """

    kwargs = _get_kwargs(
        language=language,
        q=q,
        limit=limit,
        min_status=min_status,
        variety=variety,
        release=release,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    language: str,
    q: str,
    limit: Union[Unset, int] = 20,
    min_status: Union[Unset, MinStatus] = UNSET,
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Optional[Union[EntrySearchResults, Error]]:
    """Search a language's dictionary entries

     Finds entries by headword, tolerating case, accents, spelling variants, missing
    prefixes and small typing errors (LEX-4, NFR-7). Each hit says what matched:
    the headword, a spelling variant, the root (an inflected or prefixed form) or a
    similar spelling. Best matches first. Disputed items are never returned.

    Args:
        language (str): ISO 639-3 code Example: bem.
        q (str):
        limit (Union[Unset, int]):  Default: 20.
        min_status (Union[Unset, MinStatus]):
        variety (Union[Unset, str]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[EntrySearchResults, Error]
    """

    return (
        await asyncio_detailed(
            client=client,
            language=language,
            q=q,
            limit=limit,
            min_status=min_status,
            variety=variety,
            release=release,
        )
    ).parsed
