from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response

from ...models.error import Error


def _get_kwargs(
    name: str,
) -> dict[str, Any]:
    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/v1/releases/{name}/download".format(
            name=name,
        ),
    }

    return _kwargs


def _parse_response(*, client: Union[AuthenticatedClient, Client], response: httpx.Response) -> Error:
    if response.status_code == 404:
        response_404 = Error.from_dict(response.json())

        return response_404

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Error]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    name: str,
    *,
    client: Union[AuthenticatedClient, Client],
) -> Response[Error]:
    """Download a release as open data

     The release as a zip of CSV files, one per table (entries, variants, senses,
    examples, phrases, concepts, item_sources, sources), with `release.json` and a
    README describing them (NFR-9, US-17). Frozen like a pinned read (ADR-0012):
    concepts and source details are as they are now. Public.

    Args:
        name (str):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Error]
    """

    kwargs = _get_kwargs(
        name=name,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    name: str,
    *,
    client: Union[AuthenticatedClient, Client],
) -> Optional[Error]:
    """Download a release as open data

     The release as a zip of CSV files, one per table (entries, variants, senses,
    examples, phrases, concepts, item_sources, sources), with `release.json` and a
    README describing them (NFR-9, US-17). Frozen like a pinned read (ADR-0012):
    concepts and source details are as they are now. Public.

    Args:
        name (str):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Error
    """

    return sync_detailed(
        name=name,
        client=client,
    ).parsed


async def asyncio_detailed(
    name: str,
    *,
    client: Union[AuthenticatedClient, Client],
) -> Response[Error]:
    """Download a release as open data

     The release as a zip of CSV files, one per table (entries, variants, senses,
    examples, phrases, concepts, item_sources, sources), with `release.json` and a
    README describing them (NFR-9, US-17). Frozen like a pinned read (ADR-0012):
    concepts and source details are as they are now. Public.

    Args:
        name (str):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Error]
    """

    kwargs = _get_kwargs(
        name=name,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    name: str,
    *,
    client: Union[AuthenticatedClient, Client],
) -> Optional[Error]:
    """Download a release as open data

     The release as a zip of CSV files, one per table (entries, variants, senses,
    examples, phrases, concepts, item_sources, sources), with `release.json` and a
    README describing them (NFR-9, US-17). Frozen like a pinned read (ADR-0012):
    concepts and source details are as they are now. Public.

    Args:
        name (str):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Error
    """

    return (
        await asyncio_detailed(
            name=name,
            client=client,
        )
    ).parsed
