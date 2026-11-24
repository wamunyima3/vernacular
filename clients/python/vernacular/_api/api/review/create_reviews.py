from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response

from ...models.error import Error
from ...models.review_input import ReviewInput
from ...models.review_list import ReviewList


def _get_kwargs(
    *,
    body: ReviewInput,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/v1/reviews",
    }

    _kwargs["json"] = body.to_dict()

    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, ReviewList]:
    if response.status_code == 201:
        response_201 = ReviewList.from_dict(response.json())

        return response_201

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Error, ReviewList]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    body: ReviewInput,
) -> Response[Union[Error, ReviewList]]:
    """Record review decisions (reviewers)

     Accept or reject candidates from one card. Candidates left out stay in the queue.
    Each reviewer decides on an item once (`conflict` otherwise).

    Args:
        body (ReviewInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, ReviewList]]
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
    body: ReviewInput,
) -> Optional[Union[Error, ReviewList]]:
    """Record review decisions (reviewers)

     Accept or reject candidates from one card. Candidates left out stay in the queue.
    Each reviewer decides on an item once (`conflict` otherwise).

    Args:
        body (ReviewInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, ReviewList]
    """

    return sync_detailed(
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    body: ReviewInput,
) -> Response[Union[Error, ReviewList]]:
    """Record review decisions (reviewers)

     Accept or reject candidates from one card. Candidates left out stay in the queue.
    Each reviewer decides on an item once (`conflict` otherwise).

    Args:
        body (ReviewInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, ReviewList]]
    """

    kwargs = _get_kwargs(
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    body: ReviewInput,
) -> Optional[Union[Error, ReviewList]]:
    """Record review decisions (reviewers)

     Accept or reject candidates from one card. Candidates left out stay in the queue.
    Each reviewer decides on an item once (`conflict` otherwise).

    Args:
        body (ReviewInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, ReviewList]
    """

    return (
        await asyncio_detailed(
            client=client,
            body=body,
        )
    ).parsed
