from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response

from ...models.decide_report_body import DecideReportBody
from ...models.decide_report_response_200 import DecideReportResponse200
from ...models.error import Error
from uuid import UUID


def _get_kwargs(
    id: UUID,
    *,
    body: DecideReportBody,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/v1/reports/{id}/decision".format(
            id=id,
        ),
    }

    _kwargs["json"] = body.to_dict()

    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[DecideReportResponse200, Error]:
    if response.status_code == 200:
        response_200 = DecideReportResponse200.from_dict(response.json())

        return response_200

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[DecideReportResponse200, Error]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    id: UUID,
    *,
    client: AuthenticatedClient,
    body: DecideReportBody,
) -> Response[Union[DecideReportResponse200, Error]]:
    """Confirm or dismiss a report (reviewers)

     Confirming says the report has merit: the sense becomes `disputed` and goes to the
    senior reviewers (US-15). Dismissing closes it. Reviewers cannot decide their own
    reports.

    Args:
        id (UUID):
        body (DecideReportBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[DecideReportResponse200, Error]]
    """

    kwargs = _get_kwargs(
        id=id,
        body=body,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    id: UUID,
    *,
    client: AuthenticatedClient,
    body: DecideReportBody,
) -> Optional[Union[DecideReportResponse200, Error]]:
    """Confirm or dismiss a report (reviewers)

     Confirming says the report has merit: the sense becomes `disputed` and goes to the
    senior reviewers (US-15). Dismissing closes it. Reviewers cannot decide their own
    reports.

    Args:
        id (UUID):
        body (DecideReportBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[DecideReportResponse200, Error]
    """

    return sync_detailed(
        id=id,
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    id: UUID,
    *,
    client: AuthenticatedClient,
    body: DecideReportBody,
) -> Response[Union[DecideReportResponse200, Error]]:
    """Confirm or dismiss a report (reviewers)

     Confirming says the report has merit: the sense becomes `disputed` and goes to the
    senior reviewers (US-15). Dismissing closes it. Reviewers cannot decide their own
    reports.

    Args:
        id (UUID):
        body (DecideReportBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[DecideReportResponse200, Error]]
    """

    kwargs = _get_kwargs(
        id=id,
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    id: UUID,
    *,
    client: AuthenticatedClient,
    body: DecideReportBody,
) -> Optional[Union[DecideReportResponse200, Error]]:
    """Confirm or dismiss a report (reviewers)

     Confirming says the report has merit: the sense becomes `disputed` and goes to the
    senior reviewers (US-15). Dismissing closes it. Reviewers cannot decide their own
    reports.

    Args:
        id (UUID):
        body (DecideReportBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[DecideReportResponse200, Error]
    """

    return (
        await asyncio_detailed(
            id=id,
            client=client,
            body=body,
        )
    ).parsed
