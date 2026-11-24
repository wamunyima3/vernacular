from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.error import Error
from ...models.report_list import ReportList


def _get_kwargs(
    *,
    language: str,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["language"] = language

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/v1/reports",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, ReportList]:
    if response.status_code == 200:
        response_200 = ReportList.from_dict(response.json())

        return response_200

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Error, ReportList]]:
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
) -> Response[Union[Error, ReportList]]:
    """Open reports in a language (reviewers)

     The reviewers' report inbox, oldest first, at most 100.

    Args:
        language (str): ISO 639-3 code Example: bem.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, ReportList]]
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
    client: AuthenticatedClient,
    language: str,
) -> Optional[Union[Error, ReportList]]:
    """Open reports in a language (reviewers)

     The reviewers' report inbox, oldest first, at most 100.

    Args:
        language (str): ISO 639-3 code Example: bem.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, ReportList]
    """

    return sync_detailed(
        client=client,
        language=language,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    language: str,
) -> Response[Union[Error, ReportList]]:
    """Open reports in a language (reviewers)

     The reviewers' report inbox, oldest first, at most 100.

    Args:
        language (str): ISO 639-3 code Example: bem.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, ReportList]]
    """

    kwargs = _get_kwargs(
        language=language,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    language: str,
) -> Optional[Union[Error, ReportList]]:
    """Open reports in a language (reviewers)

     The reviewers' report inbox, oldest first, at most 100.

    Args:
        language (str): ISO 639-3 code Example: bem.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, ReportList]
    """

    return (
        await asyncio_detailed(
            client=client,
            language=language,
        )
    ).parsed
