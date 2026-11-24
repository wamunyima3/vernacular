from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.concept import Concept
from ...models.error import Error
from ...models.min_status import MinStatus
from ...types import Unset
from uuid import UUID


def _get_kwargs(
    id: UUID,
    *,
    min_status: Union[Unset, MinStatus] = UNSET,
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    json_min_status: Union[Unset, str] = UNSET
    if not isinstance(min_status, Unset):
        json_min_status = min_status.value

    params["min_status"] = json_min_status

    params["variety"] = variety

    params["release"] = release

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/v1/concepts/{id}".format(
            id=id,
        ),
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Concept, Error]:
    if response.status_code == 200:
        response_200 = Concept.from_dict(response.json())

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
) -> Response[Union[Concept, Error]]:
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
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Response[Union[Concept, Error]]:
    """A meaning across all languages

     The language-independent concept, and every servable entry in any language with a
    sense that expresses it (LEX-7, US-04). Each entry lists only those senses.

    Args:
        id (UUID):
        min_status (Union[Unset, MinStatus]):
        variety (Union[Unset, str]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Concept, Error]]
    """

    kwargs = _get_kwargs(
        id=id,
        min_status=min_status,
        variety=variety,
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
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Optional[Union[Concept, Error]]:
    """A meaning across all languages

     The language-independent concept, and every servable entry in any language with a
    sense that expresses it (LEX-7, US-04). Each entry lists only those senses.

    Args:
        id (UUID):
        min_status (Union[Unset, MinStatus]):
        variety (Union[Unset, str]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Concept, Error]
    """

    return sync_detailed(
        id=id,
        client=client,
        min_status=min_status,
        variety=variety,
        release=release,
    ).parsed


async def asyncio_detailed(
    id: UUID,
    *,
    client: AuthenticatedClient,
    min_status: Union[Unset, MinStatus] = UNSET,
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Response[Union[Concept, Error]]:
    """A meaning across all languages

     The language-independent concept, and every servable entry in any language with a
    sense that expresses it (LEX-7, US-04). Each entry lists only those senses.

    Args:
        id (UUID):
        min_status (Union[Unset, MinStatus]):
        variety (Union[Unset, str]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Concept, Error]]
    """

    kwargs = _get_kwargs(
        id=id,
        min_status=min_status,
        variety=variety,
        release=release,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    id: UUID,
    *,
    client: AuthenticatedClient,
    min_status: Union[Unset, MinStatus] = UNSET,
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Optional[Union[Concept, Error]]:
    """A meaning across all languages

     The language-independent concept, and every servable entry in any language with a
    sense that expresses it (LEX-7, US-04). Each entry lists only those senses.

    Args:
        id (UUID):
        min_status (Union[Unset, MinStatus]):
        variety (Union[Unset, str]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Concept, Error]
    """

    return (
        await asyncio_detailed(
            id=id,
            client=client,
            min_status=min_status,
            variety=variety,
            release=release,
        )
    ).parsed
