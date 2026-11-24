from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response

from ...models.error import Error
from ...models.qualification_test import QualificationTest
from ...models.start_qualification_body import StartQualificationBody


def _get_kwargs(
    *,
    body: StartQualificationBody,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/v1/qualifications",
    }

    _kwargs["json"] = body.to_dict()

    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, QualificationTest]:
    if response.status_code == 201:
        response_201 = QualificationTest.from_dict(response.json())

        return response_201

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Error, QualificationTest]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    body: StartQualificationBody,
) -> Response[Union[Error, QualificationTest]]:
    """Start a reviewer qualification test

     Ten questions drawn from the language's verified words: pick the word for each English
    meaning (DQ-3). A pass is 8 or more; a senior reviewer then confirms it, which approves
    you to review. After a fail you can try again the next day.

    Args:
        body (StartQualificationBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, QualificationTest]]
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
    body: StartQualificationBody,
) -> Optional[Union[Error, QualificationTest]]:
    """Start a reviewer qualification test

     Ten questions drawn from the language's verified words: pick the word for each English
    meaning (DQ-3). A pass is 8 or more; a senior reviewer then confirms it, which approves
    you to review. After a fail you can try again the next day.

    Args:
        body (StartQualificationBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, QualificationTest]
    """

    return sync_detailed(
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    body: StartQualificationBody,
) -> Response[Union[Error, QualificationTest]]:
    """Start a reviewer qualification test

     Ten questions drawn from the language's verified words: pick the word for each English
    meaning (DQ-3). A pass is 8 or more; a senior reviewer then confirms it, which approves
    you to review. After a fail you can try again the next day.

    Args:
        body (StartQualificationBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, QualificationTest]]
    """

    kwargs = _get_kwargs(
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    body: StartQualificationBody,
) -> Optional[Union[Error, QualificationTest]]:
    """Start a reviewer qualification test

     Ten questions drawn from the language's verified words: pick the word for each English
    meaning (DQ-3). A pass is 8 or more; a senior reviewer then confirms it, which approves
    you to review. After a fail you can try again the next day.

    Args:
        body (StartQualificationBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, QualificationTest]
    """

    return (
        await asyncio_detailed(
            client=client,
            body=body,
        )
    ).parsed
