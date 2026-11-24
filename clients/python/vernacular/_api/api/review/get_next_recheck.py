from http import HTTPStatus
from typing import Any, Optional, Union, cast

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.error import Error
from ...models.recheck import Recheck
from ...types import Unset
from uuid import UUID


def _get_kwargs(
    *,
    language: str,
    skip: Union[Unset, list[UUID]] = UNSET,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["language"] = language

    json_skip: Union[Unset, list[str]] = UNSET
    if not isinstance(skip, Unset):
        json_skip = []
        for skip_item_data in skip:
            skip_item = str(skip_item_data)
            json_skip.append(skip_item)

    params["skip"] = json_skip

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/v1/rechecks/next",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Any, Error, Recheck]:
    if response.status_code == 200:
        response_200 = Recheck.from_dict(response.json())

        return response_200

    if response.status_code == 204:
        response_204 = cast(Any, None)
        return response_204

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Any, Error, Recheck]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    language: str,
    skip: Union[Unset, list[UUID]] = UNSET,
) -> Response[Union[Any, Error, Recheck]]:
    """The next item to re-check (senior reviewers)

     The oldest open re-check in the language, never one the caller wrote or reviewed.
    `204` when none is left.

    Args:
        language (str): ISO 639-3 code Example: bem.
        skip (Union[Unset, list[UUID]]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Any, Error, Recheck]]
    """

    kwargs = _get_kwargs(
        language=language,
        skip=skip,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: AuthenticatedClient,
    language: str,
    skip: Union[Unset, list[UUID]] = UNSET,
) -> Optional[Union[Any, Error, Recheck]]:
    """The next item to re-check (senior reviewers)

     The oldest open re-check in the language, never one the caller wrote or reviewed.
    `204` when none is left.

    Args:
        language (str): ISO 639-3 code Example: bem.
        skip (Union[Unset, list[UUID]]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Any, Error, Recheck]
    """

    return sync_detailed(
        client=client,
        language=language,
        skip=skip,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    language: str,
    skip: Union[Unset, list[UUID]] = UNSET,
) -> Response[Union[Any, Error, Recheck]]:
    """The next item to re-check (senior reviewers)

     The oldest open re-check in the language, never one the caller wrote or reviewed.
    `204` when none is left.

    Args:
        language (str): ISO 639-3 code Example: bem.
        skip (Union[Unset, list[UUID]]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Any, Error, Recheck]]
    """

    kwargs = _get_kwargs(
        language=language,
        skip=skip,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    language: str,
    skip: Union[Unset, list[UUID]] = UNSET,
) -> Optional[Union[Any, Error, Recheck]]:
    """The next item to re-check (senior reviewers)

     The oldest open re-check in the language, never one the caller wrote or reviewed.
    `204` when none is left.

    Args:
        language (str): ISO 639-3 code Example: bem.
        skip (Union[Unset, list[UUID]]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Any, Error, Recheck]
    """

    return (
        await asyncio_detailed(
            client=client,
            language=language,
            skip=skip,
        )
    ).parsed
