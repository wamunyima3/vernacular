from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response

from ...models.error import Error
from ...models.resolution import Resolution
from ...models.resolution_input import ResolutionInput
from uuid import UUID


def _get_kwargs(
    item_id: UUID,
    *,
    body: ResolutionInput,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/v1/disputes/{item_id}/resolution".format(
            item_id=item_id,
        ),
    }

    _kwargs["json"] = body.to_dict()

    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, Resolution]:
    if response.status_code == 200:
        response_200 = Resolution.from_dict(response.json())

        return response_200

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Error, Resolution]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    item_id: UUID,
    *,
    client: AuthenticatedClient,
    body: ResolutionInput,
) -> Response[Union[Error, Resolution]]:
    """Settle a dispute (senior reviewers)

     `keep` verifies the item. `send_back` returns it to the review queue, where only
    reviews made after this count. Seniors cannot settle disputes they reviewed.

    Args:
        item_id (UUID):
        body (ResolutionInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, Resolution]]
    """

    kwargs = _get_kwargs(
        item_id=item_id,
        body=body,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    item_id: UUID,
    *,
    client: AuthenticatedClient,
    body: ResolutionInput,
) -> Optional[Union[Error, Resolution]]:
    """Settle a dispute (senior reviewers)

     `keep` verifies the item. `send_back` returns it to the review queue, where only
    reviews made after this count. Seniors cannot settle disputes they reviewed.

    Args:
        item_id (UUID):
        body (ResolutionInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, Resolution]
    """

    return sync_detailed(
        item_id=item_id,
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    item_id: UUID,
    *,
    client: AuthenticatedClient,
    body: ResolutionInput,
) -> Response[Union[Error, Resolution]]:
    """Settle a dispute (senior reviewers)

     `keep` verifies the item. `send_back` returns it to the review queue, where only
    reviews made after this count. Seniors cannot settle disputes they reviewed.

    Args:
        item_id (UUID):
        body (ResolutionInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, Resolution]]
    """

    kwargs = _get_kwargs(
        item_id=item_id,
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    item_id: UUID,
    *,
    client: AuthenticatedClient,
    body: ResolutionInput,
) -> Optional[Union[Error, Resolution]]:
    """Settle a dispute (senior reviewers)

     `keep` verifies the item. `send_back` returns it to the review queue, where only
    reviews made after this count. Seniors cannot settle disputes they reviewed.

    Args:
        item_id (UUID):
        body (ResolutionInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, Resolution]
    """

    return (
        await asyncio_detailed(
            item_id=item_id,
            client=client,
            body=body,
        )
    ).parsed
