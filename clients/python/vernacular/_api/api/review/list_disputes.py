from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response

from ...models.dispute_list import DisputeList
from ...models.error import Error


def _get_kwargs() -> dict[str, Any]:
    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/v1/disputes",
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[DisputeList, Error]:
    if response.status_code == 200:
        response_200 = DisputeList.from_dict(response.json())

        return response_200

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[DisputeList, Error]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
) -> Response[Union[DisputeList, Error]]:
    """Disagreements waiting for a senior reviewer (senior reviewers)

     Items reviewers disagree on, in the languages the caller is a senior reviewer for,
    oldest first, with every decision that counts (CON-4). Resolve within 14 days (DQ-8).

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[DisputeList, Error]]
    """

    kwargs = _get_kwargs()

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: AuthenticatedClient,
) -> Optional[Union[DisputeList, Error]]:
    """Disagreements waiting for a senior reviewer (senior reviewers)

     Items reviewers disagree on, in the languages the caller is a senior reviewer for,
    oldest first, with every decision that counts (CON-4). Resolve within 14 days (DQ-8).

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[DisputeList, Error]
    """

    return sync_detailed(
        client=client,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
) -> Response[Union[DisputeList, Error]]:
    """Disagreements waiting for a senior reviewer (senior reviewers)

     Items reviewers disagree on, in the languages the caller is a senior reviewer for,
    oldest first, with every decision that counts (CON-4). Resolve within 14 days (DQ-8).

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[DisputeList, Error]]
    """

    kwargs = _get_kwargs()

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
) -> Optional[Union[DisputeList, Error]]:
    """Disagreements waiting for a senior reviewer (senior reviewers)

     Items reviewers disagree on, in the languages the caller is a senior reviewer for,
    oldest first, with every decision that counts (CON-4). Resolve within 14 days (DQ-8).

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[DisputeList, Error]
    """

    return (
        await asyncio_detailed(
            client=client,
        )
    ).parsed
