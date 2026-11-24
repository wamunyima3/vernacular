from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.clip import Clip
from ...models.error import Error
from ...types import File
from uuid import UUID


def _get_kwargs(
    *,
    body: File,
    entry_id: UUID,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    params: dict[str, Any] = {}

    json_entry_id = str(entry_id)
    params["entry_id"] = json_entry_id

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/v1/audio",
        "params": params,
    }

    _kwargs["content"] = body.payload

    headers["Content-Type"] = "application/octet-stream"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Clip, Error]:
    if response.status_code == 201:
        response_201 = Clip.from_dict(response.json())

        return response_201

    if response.status_code == 403:
        response_403 = Error.from_dict(response.json())

        return response_403

    if response.status_code == 404:
        response_404 = Error.from_dict(response.json())

        return response_404

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Clip, Error]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    body: File,
    entry_id: UUID,
) -> Response[Union[Clip, Error]]:
    """Record a pronunciation of an entry

     The recording as the body, at most 1 MB: Ogg, WebM, MP4, MP3 or WAV. You need a consent
    form on file that names the entry's language. Recording an entry again replaces your
    earlier clip.

    Args:
        entry_id (UUID):
        body (File):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Clip, Error]]
    """

    kwargs = _get_kwargs(
        body=body,
        entry_id=entry_id,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: AuthenticatedClient,
    body: File,
    entry_id: UUID,
) -> Optional[Union[Clip, Error]]:
    """Record a pronunciation of an entry

     The recording as the body, at most 1 MB: Ogg, WebM, MP4, MP3 or WAV. You need a consent
    form on file that names the entry's language. Recording an entry again replaces your
    earlier clip.

    Args:
        entry_id (UUID):
        body (File):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Clip, Error]
    """

    return sync_detailed(
        client=client,
        body=body,
        entry_id=entry_id,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    body: File,
    entry_id: UUID,
) -> Response[Union[Clip, Error]]:
    """Record a pronunciation of an entry

     The recording as the body, at most 1 MB: Ogg, WebM, MP4, MP3 or WAV. You need a consent
    form on file that names the entry's language. Recording an entry again replaces your
    earlier clip.

    Args:
        entry_id (UUID):
        body (File):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Clip, Error]]
    """

    kwargs = _get_kwargs(
        body=body,
        entry_id=entry_id,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    body: File,
    entry_id: UUID,
) -> Optional[Union[Clip, Error]]:
    """Record a pronunciation of an entry

     The recording as the body, at most 1 MB: Ogg, WebM, MP4, MP3 or WAV. You need a consent
    form on file that names the entry's language. Recording an entry again replaces your
    earlier clip.

    Args:
        entry_id (UUID):
        body (File):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Clip, Error]
    """

    return (
        await asyncio_detailed(
            client=client,
            body=body,
            entry_id=entry_id,
        )
    ).parsed
