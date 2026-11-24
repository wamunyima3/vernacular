from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response

from ...models.contribution import Contribution
from ...models.contribution_input import ContributionInput
from ...models.error import Error


def _get_kwargs(
    *,
    body: ContributionInput,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/v1/contributions",
    }

    _kwargs["json"] = body.to_dict()

    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Contribution, Error]:
    if response.status_code == 201:
        response_201 = Contribution.from_dict(response.json())

        return response_201

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Contribution, Error]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    body: ContributionInput,
) -> Response[Union[Contribution, Error]]:
    """Propose a word (contributors)

     Adds an entry and sense, and optionally an example, as `sourced` and queues them for
    blind review (CON-1, US-18). If the word and meaning already exist, the contribution
    is recorded as agreement with them instead.

    Args:
        body (ContributionInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Contribution, Error]]
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
    body: ContributionInput,
) -> Optional[Union[Contribution, Error]]:
    """Propose a word (contributors)

     Adds an entry and sense, and optionally an example, as `sourced` and queues them for
    blind review (CON-1, US-18). If the word and meaning already exist, the contribution
    is recorded as agreement with them instead.

    Args:
        body (ContributionInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Contribution, Error]
    """

    return sync_detailed(
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    body: ContributionInput,
) -> Response[Union[Contribution, Error]]:
    """Propose a word (contributors)

     Adds an entry and sense, and optionally an example, as `sourced` and queues them for
    blind review (CON-1, US-18). If the word and meaning already exist, the contribution
    is recorded as agreement with them instead.

    Args:
        body (ContributionInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Contribution, Error]]
    """

    kwargs = _get_kwargs(
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    body: ContributionInput,
) -> Optional[Union[Contribution, Error]]:
    """Propose a word (contributors)

     Adds an entry and sense, and optionally an example, as `sourced` and queues them for
    blind review (CON-1, US-18). If the word and meaning already exist, the contribution
    is recorded as agreement with them instead.

    Args:
        body (ContributionInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Contribution, Error]
    """

    return (
        await asyncio_detailed(
            client=client,
            body=body,
        )
    ).parsed
