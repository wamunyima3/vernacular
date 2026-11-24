from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.clip_list import ClipList
from ...models.error import Error
from uuid import UUID


def _get_kwargs(
    *,
    entry_id: UUID,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    json_entry_id = str(entry_id)
    params["entry_id"] = json_entry_id

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/v1/audio",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[ClipList, Error]:
    if response.status_code == 200:
        response_200 = ClipList.from_dict(response.json())

        return response_200

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[ClipList, Error]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    entry_id: UUID,
) -> Response[Union[ClipList, Error]]:
    """An entry's pronunciation clips

    Args:
        entry_id (UUID):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[ClipList, Error]]
    """

    kwargs = _get_kwargs(
        entry_id=entry_id,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: Union[AuthenticatedClient, Client],
    entry_id: UUID,
) -> Optional[Union[ClipList, Error]]:
    """An entry's pronunciation clips

    Args:
        entry_id (UUID):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[ClipList, Error]
    """

    return sync_detailed(
        client=client,
        entry_id=entry_id,
    ).parsed


async def asyncio_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    entry_id: UUID,
) -> Response[Union[ClipList, Error]]:
    """An entry's pronunciation clips

    Args:
        entry_id (UUID):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[ClipList, Error]]
    """

    kwargs = _get_kwargs(
        entry_id=entry_id,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: Union[AuthenticatedClient, Client],
    entry_id: UUID,
) -> Optional[Union[ClipList, Error]]:
    """An entry's pronunciation clips

    Args:
        entry_id (UUID):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[ClipList, Error]
    """

    return (
        await asyncio_detailed(
            client=client,
            entry_id=entry_id,
        )
    ).parsed
