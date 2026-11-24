from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.entry import Entry
from ...models.error import Error
from ...models.min_status import MinStatus
from ...types import Unset
from uuid import UUID


def _get_kwargs(
    id: UUID,
    *,
    min_status: Union[Unset, MinStatus] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    json_min_status: Union[Unset, str] = UNSET
    if not isinstance(min_status, Unset):
        json_min_status = min_status.value

    params["min_status"] = json_min_status

    params["release"] = release

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/v1/entries/{id}".format(
            id=id,
        ),
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Entry, Error]:
    if response.status_code == 200:
        response_200 = Entry.from_dict(response.json())

        return response_200

    if response.status_code == 401:
        response_401 = Error.from_dict(response.json())

        return response_401

    if response.status_code == 404:
        response_404 = Error.from_dict(response.json())

        return response_404

    if response.status_code == 429:
        response_429 = Error.from_dict(response.json())

        return response_429

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Entry, Error]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    id: UUID,
    *,
    client: AuthenticatedClient,
    min_status: Union[Unset, MinStatus] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Response[Union[Entry, Error]]:
    """A full dictionary entry

     The entry with its spelling variants, and each servable sense with its definition,
    concept and example sentences (LEX-1, LEX-3, LEX-5, US-16). Every item carries its
    trust fields. A disputed entry is not found.

    Args:
        id (UUID):
        min_status (Union[Unset, MinStatus]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Entry, Error]]
    """

    kwargs = _get_kwargs(
        id=id,
        min_status=min_status,
        release=release,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    id: UUID,
    *,
    client: AuthenticatedClient,
    min_status: Union[Unset, MinStatus] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Optional[Union[Entry, Error]]:
    """A full dictionary entry

     The entry with its spelling variants, and each servable sense with its definition,
    concept and example sentences (LEX-1, LEX-3, LEX-5, US-16). Every item carries its
    trust fields. A disputed entry is not found.

    Args:
        id (UUID):
        min_status (Union[Unset, MinStatus]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Entry, Error]
    """

    return sync_detailed(
        id=id,
        client=client,
        min_status=min_status,
        release=release,
    ).parsed


async def asyncio_detailed(
    id: UUID,
    *,
    client: AuthenticatedClient,
    min_status: Union[Unset, MinStatus] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Response[Union[Entry, Error]]:
    """A full dictionary entry

     The entry with its spelling variants, and each servable sense with its definition,
    concept and example sentences (LEX-1, LEX-3, LEX-5, US-16). Every item carries its
    trust fields. A disputed entry is not found.

    Args:
        id (UUID):
        min_status (Union[Unset, MinStatus]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Entry, Error]]
    """

    kwargs = _get_kwargs(
        id=id,
        min_status=min_status,
        release=release,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    id: UUID,
    *,
    client: AuthenticatedClient,
    min_status: Union[Unset, MinStatus] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Optional[Union[Entry, Error]]:
    """A full dictionary entry

     The entry with its spelling variants, and each servable sense with its definition,
    concept and example sentences (LEX-1, LEX-3, LEX-5, US-16). Every item carries its
    trust fields. A disputed entry is not found.

    Args:
        id (UUID):
        min_status (Union[Unset, MinStatus]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Entry, Error]
    """

    return (
        await asyncio_detailed(
            id=id,
            client=client,
            min_status=min_status,
            release=release,
        )
    ).parsed
