from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.error import Error
from ...models.key_usage import KeyUsage
from ...types import Unset
from uuid import UUID


def _get_kwargs(
    id: UUID,
    *,
    days: Union[Unset, int] = 30,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["days"] = days

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/v1/keys/{id}/usage".format(
            id=id,
        ),
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, KeyUsage]:
    if response.status_code == 200:
        response_200 = KeyUsage.from_dict(response.json())

        return response_200

    if response.status_code == 404:
        response_404 = Error.from_dict(response.json())

        return response_404

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Error, KeyUsage]]:
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
    days: Union[Unset, int] = 30,
) -> Response[Union[Error, KeyUsage]]:
    """Daily calls with one of your API keys

     Calls, errors and calls refused by the rate limit, per UTC day, newest first; days
    without calls are left out (DEV-5, US-13). Counts arrive within a minute.

    Args:
        id (UUID):
        days (Union[Unset, int]):  Default: 30.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, KeyUsage]]
    """

    kwargs = _get_kwargs(
        id=id,
        days=days,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    id: UUID,
    *,
    client: AuthenticatedClient,
    days: Union[Unset, int] = 30,
) -> Optional[Union[Error, KeyUsage]]:
    """Daily calls with one of your API keys

     Calls, errors and calls refused by the rate limit, per UTC day, newest first; days
    without calls are left out (DEV-5, US-13). Counts arrive within a minute.

    Args:
        id (UUID):
        days (Union[Unset, int]):  Default: 30.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, KeyUsage]
    """

    return sync_detailed(
        id=id,
        client=client,
        days=days,
    ).parsed


async def asyncio_detailed(
    id: UUID,
    *,
    client: AuthenticatedClient,
    days: Union[Unset, int] = 30,
) -> Response[Union[Error, KeyUsage]]:
    """Daily calls with one of your API keys

     Calls, errors and calls refused by the rate limit, per UTC day, newest first; days
    without calls are left out (DEV-5, US-13). Counts arrive within a minute.

    Args:
        id (UUID):
        days (Union[Unset, int]):  Default: 30.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, KeyUsage]]
    """

    kwargs = _get_kwargs(
        id=id,
        days=days,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    id: UUID,
    *,
    client: AuthenticatedClient,
    days: Union[Unset, int] = 30,
) -> Optional[Union[Error, KeyUsage]]:
    """Daily calls with one of your API keys

     Calls, errors and calls refused by the rate limit, per UTC day, newest first; days
    without calls are left out (DEV-5, US-13). Counts arrive within a minute.

    Args:
        id (UUID):
        days (Union[Unset, int]):  Default: 30.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, KeyUsage]
    """

    return (
        await asyncio_detailed(
            id=id,
            client=client,
            days=days,
        )
    ).parsed
