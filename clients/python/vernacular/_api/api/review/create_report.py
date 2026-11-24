from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response

from ...models.create_report_body import CreateReportBody
from ...models.create_report_response_201 import CreateReportResponse201
from ...models.error import Error


def _get_kwargs(
    *,
    body: CreateReportBody,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/v1/reports",
    }

    _kwargs["json"] = body.to_dict()

    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[CreateReportResponse201, Error]:
    if response.status_code == 201:
        response_201 = CreateReportResponse201.from_dict(response.json())

        return response_201

    if response.status_code == 404:
        response_404 = Error.from_dict(response.json())

        return response_404

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[CreateReportResponse201, Error]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    body: CreateReportBody,
) -> Response[Union[CreateReportResponse201, Error]]:
    """Report a wrong meaning

     Anyone can report a sense as wrong, signed in or not (PRV-3, US-15). The report
    reaches the reviewers of its language. If a reviewer confirms it, the sense becomes
    `disputed` until a senior reviewer settles it.

    Args:
        body (CreateReportBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[CreateReportResponse201, Error]]
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
    body: CreateReportBody,
) -> Optional[Union[CreateReportResponse201, Error]]:
    """Report a wrong meaning

     Anyone can report a sense as wrong, signed in or not (PRV-3, US-15). The report
    reaches the reviewers of its language. If a reviewer confirms it, the sense becomes
    `disputed` until a senior reviewer settles it.

    Args:
        body (CreateReportBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[CreateReportResponse201, Error]
    """

    return sync_detailed(
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    body: CreateReportBody,
) -> Response[Union[CreateReportResponse201, Error]]:
    """Report a wrong meaning

     Anyone can report a sense as wrong, signed in or not (PRV-3, US-15). The report
    reaches the reviewers of its language. If a reviewer confirms it, the sense becomes
    `disputed` until a senior reviewer settles it.

    Args:
        body (CreateReportBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[CreateReportResponse201, Error]]
    """

    kwargs = _get_kwargs(
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    body: CreateReportBody,
) -> Optional[Union[CreateReportResponse201, Error]]:
    """Report a wrong meaning

     Anyone can report a sense as wrong, signed in or not (PRV-3, US-15). The report
    reaches the reviewers of its language. If a reviewer confirms it, the sense becomes
    `disputed` until a senior reviewer settles it.

    Args:
        body (CreateReportBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[CreateReportResponse201, Error]
    """

    return (
        await asyncio_detailed(
            client=client,
            body=body,
        )
    ).parsed
