from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response

from ...models.error import Error
from ...models.recheck_input import RecheckInput
from ...models.recheck_result import RecheckResult
from uuid import UUID


def _get_kwargs(
    item_id: UUID,
    *,
    body: RecheckInput,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/v1/rechecks/{item_id}".format(
            item_id=item_id,
        ),
    }

    _kwargs["json"] = body.to_dict()

    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, RecheckResult]:
    if response.status_code == 200:
        response_200 = RecheckResult.from_dict(response.json())

        return response_200

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Error, RecheckResult]]:
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
    body: RecheckInput,
) -> Response[Union[Error, RecheckResult]]:
    """Record a re-check (senior reviewers)

     `pass` leaves the item verified. `fail`, with a note saying what is wrong, sends it
    back to the review queue, where only reviews made after this count.

    Args:
        item_id (UUID):
        body (RecheckInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, RecheckResult]]
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
    body: RecheckInput,
) -> Optional[Union[Error, RecheckResult]]:
    """Record a re-check (senior reviewers)

     `pass` leaves the item verified. `fail`, with a note saying what is wrong, sends it
    back to the review queue, where only reviews made after this count.

    Args:
        item_id (UUID):
        body (RecheckInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, RecheckResult]
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
    body: RecheckInput,
) -> Response[Union[Error, RecheckResult]]:
    """Record a re-check (senior reviewers)

     `pass` leaves the item verified. `fail`, with a note saying what is wrong, sends it
    back to the review queue, where only reviews made after this count.

    Args:
        item_id (UUID):
        body (RecheckInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, RecheckResult]]
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
    body: RecheckInput,
) -> Optional[Union[Error, RecheckResult]]:
    """Record a re-check (senior reviewers)

     `pass` leaves the item verified. `fail`, with a note saying what is wrong, sends it
    back to the review queue, where only reviews made after this count.

    Args:
        item_id (UUID):
        body (RecheckInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, RecheckResult]
    """

    return (
        await asyncio_detailed(
            item_id=item_id,
            client=client,
            body=body,
        )
    ).parsed
