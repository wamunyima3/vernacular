from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.error import Error
from ...models.example_coverage import ExampleCoverage


def _get_kwargs(
    *,
    language: str,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["language"] = language

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/v1/examples/missing",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, ExampleCoverage]:
    if response.status_code == 200:
        response_200 = ExampleCoverage.from_dict(response.json())

        return response_200

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Error, ExampleCoverage]]:
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
) -> Response[Union[Error, ExampleCoverage]]:
    """Core meanings without a verified example

     Coverage of LEX-5 in a language: verified meanings of core-list concepts, how many
    have a verified example, and those that do not yet, for contributors to work on.

    Args:
        language (str): ISO 639-3 code Example: bem.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, ExampleCoverage]]
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
) -> Optional[Union[Error, ExampleCoverage]]:
    """Core meanings without a verified example

     Coverage of LEX-5 in a language: verified meanings of core-list concepts, how many
    have a verified example, and those that do not yet, for contributors to work on.

    Args:
        language (str): ISO 639-3 code Example: bem.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, ExampleCoverage]
    """

    return sync_detailed(
        client=client,
        language=language,
    ).parsed


async def asyncio_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    language: str,
) -> Response[Union[Error, ExampleCoverage]]:
    """Core meanings without a verified example

     Coverage of LEX-5 in a language: verified meanings of core-list concepts, how many
    have a verified example, and those that do not yet, for contributors to work on.

    Args:
        language (str): ISO 639-3 code Example: bem.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, ExampleCoverage]]
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
) -> Optional[Union[Error, ExampleCoverage]]:
    """Core meanings without a verified example

     Coverage of LEX-5 in a language: verified meanings of core-list concepts, how many
    have a verified example, and those that do not yet, for contributors to work on.

    Args:
        language (str): ISO 639-3 code Example: bem.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, ExampleCoverage]
    """

    return (
        await asyncio_detailed(
            client=client,
            language=language,
        )
    ).parsed
