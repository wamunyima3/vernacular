from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response

from ...models.error import Error
from ...models.release import Release
from ...models.release_input import ReleaseInput


def _get_kwargs(
    *,
    body: ReleaseInput,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/v1/releases",
    }

    _kwargs["json"] = body.to_dict()

    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, Release]:
    if response.status_code == 201:
        response_201 = Release.from_dict(response.json())

        return response_201

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Error, Release]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    body: ReleaseInput,
) -> Response[Union[Error, Release]]:
    """Build a release from the verified items (admin)

     Freezes every verified item into a new release named `<year>.<n>` (US-29). Refused
    with `conflict` when nothing is verified yet or nothing changed since the last release.

    Args:
        body (ReleaseInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, Release]]
    """

    kwargs = _get_kwargs(
        body=body,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: AuthenticatedClient,
    body: ReleaseInput,
) -> Optional[Union[Error, Release]]:
    """Build a release from the verified items (admin)

     Freezes every verified item into a new release named `<year>.<n>` (US-29). Refused
    with `conflict` when nothing is verified yet or nothing changed since the last release.

    Args:
        body (ReleaseInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, Release]
    """

    return sync_detailed(
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    body: ReleaseInput,
) -> Response[Union[Error, Release]]:
    """Build a release from the verified items (admin)

     Freezes every verified item into a new release named `<year>.<n>` (US-29). Refused
    with `conflict` when nothing is verified yet or nothing changed since the last release.

    Args:
        body (ReleaseInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, Release]]
    """

    kwargs = _get_kwargs(
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    body: ReleaseInput,
) -> Optional[Union[Error, Release]]:
    """Build a release from the verified items (admin)

     Freezes every verified item into a new release named `<year>.<n>` (US-29). Refused
    with `conflict` when nothing is verified yet or nothing changed since the last release.

    Args:
        body (ReleaseInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, Release]
    """

    return (
        await asyncio_detailed(
            client=client,
            body=body,
        )
    ).parsed
