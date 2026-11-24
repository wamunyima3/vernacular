from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.credit_list import CreditList
from ...models.error import Error


def _get_kwargs(
    *,
    language: str,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["language"] = language

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/v1/credits",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[CreditList, Error]:
    if response.status_code == 200:
        response_200 = CreditList.from_dict(response.json())

        return response_200

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[CreditList, Error]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    language: str,
) -> Response[Union[CreditList, Error]]:
    """The volunteers behind a language's data

     Everyone credited for accepted work in the language (CON-5): contributors for their
    words that were verified, reviewers for decisions that matched a settled outcome.
    Only people who chose to be credited, by display name (NFR-6).

    Args:
        language (str): ISO 639-3 code Example: bem.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[CreditList, Error]]
    """

    kwargs = _get_kwargs(
        language=language,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: Union[AuthenticatedClient, Client],
    language: str,
) -> Optional[Union[CreditList, Error]]:
    """The volunteers behind a language's data

     Everyone credited for accepted work in the language (CON-5): contributors for their
    words that were verified, reviewers for decisions that matched a settled outcome.
    Only people who chose to be credited, by display name (NFR-6).

    Args:
        language (str): ISO 639-3 code Example: bem.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[CreditList, Error]
    """

    return sync_detailed(
        client=client,
        language=language,
    ).parsed


async def asyncio_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    language: str,
) -> Response[Union[CreditList, Error]]:
    """The volunteers behind a language's data

     Everyone credited for accepted work in the language (CON-5): contributors for their
    words that were verified, reviewers for decisions that matched a settled outcome.
    Only people who chose to be credited, by display name (NFR-6).

    Args:
        language (str): ISO 639-3 code Example: bem.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[CreditList, Error]]
    """

    kwargs = _get_kwargs(
        language=language,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: Union[AuthenticatedClient, Client],
    language: str,
) -> Optional[Union[CreditList, Error]]:
    """The volunteers behind a language's data

     Everyone credited for accepted work in the language (CON-5): contributors for their
    words that were verified, reviewers for decisions that matched a settled outcome.
    Only people who chose to be credited, by display name (NFR-6).

    Args:
        language (str): ISO 639-3 code Example: bem.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[CreditList, Error]
    """

    return (
        await asyncio_detailed(
            client=client,
            language=language,
        )
    ).parsed
