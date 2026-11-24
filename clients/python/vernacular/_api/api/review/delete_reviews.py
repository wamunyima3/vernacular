from http import HTTPStatus
from typing import Any, Optional, Union, cast

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.error import Error
from uuid import UUID


def _get_kwargs(
    *,
    ids: list[UUID],
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    json_ids = []
    for ids_item_data in ids:
        ids_item = str(ids_item_data)
        json_ids.append(ids_item)

    params["ids"] = json_ids

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "delete",
        "url": "/v1/reviews",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Any, Error]:
    if response.status_code == 204:
        response_204 = cast(Any, None)
        return response_204

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Any, Error]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    ids: list[UUID],
) -> Response[Union[Any, Error]]:
    """Undo your own decision (reviewers)

     Takes back the reviews one decision recorded, all or none. Only their reviewer can,
    within 10 minutes and while the items are still in the queue (`conflict` otherwise).

    Args:
        ids (list[UUID]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Any, Error]]
    """

    kwargs = _get_kwargs(
        ids=ids,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: AuthenticatedClient,
    ids: list[UUID],
) -> Optional[Union[Any, Error]]:
    """Undo your own decision (reviewers)

     Takes back the reviews one decision recorded, all or none. Only their reviewer can,
    within 10 minutes and while the items are still in the queue (`conflict` otherwise).

    Args:
        ids (list[UUID]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Any, Error]
    """

    return sync_detailed(
        client=client,
        ids=ids,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    ids: list[UUID],
) -> Response[Union[Any, Error]]:
    """Undo your own decision (reviewers)

     Takes back the reviews one decision recorded, all or none. Only their reviewer can,
    within 10 minutes and while the items are still in the queue (`conflict` otherwise).

    Args:
        ids (list[UUID]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Any, Error]]
    """

    kwargs = _get_kwargs(
        ids=ids,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    ids: list[UUID],
) -> Optional[Union[Any, Error]]:
    """Undo your own decision (reviewers)

     Takes back the reviews one decision recorded, all or none. Only their reviewer can,
    within 10 minutes and while the items are still in the queue (`conflict` otherwise).

    Args:
        ids (list[UUID]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Any, Error]
    """

    return (
        await asyncio_detailed(
            client=client,
            ids=ids,
        )
    ).parsed
