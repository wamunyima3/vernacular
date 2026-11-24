from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response

from ...models.error import Error
from ...models.example_contribution import ExampleContribution
from ...models.example_contribution_input import ExampleContributionInput
from uuid import UUID


def _get_kwargs(
    id: UUID,
    *,
    body: ExampleContributionInput,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/v1/senses/{id}/examples".format(
            id=id,
        ),
    }

    _kwargs["json"] = body.to_dict()

    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, ExampleContribution]:
    if response.status_code == 201:
        response_201 = ExampleContribution.from_dict(response.json())

        return response_201

    if response.status_code == 404:
        response_404 = Error.from_dict(response.json())

        return response_404

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Error, ExampleContribution]]:
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
    body: ExampleContributionInput,
) -> Response[Union[Error, ExampleContribution]]:
    """Add an example sentence to a meaning (contributors)

     Adds the sentence and its English to an existing meaning as `sourced` and queues it
    for blind review (LEX-5, CON-1).

    Args:
        id (UUID):
        body (ExampleContributionInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, ExampleContribution]]
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
    body: ExampleContributionInput,
) -> Optional[Union[Error, ExampleContribution]]:
    """Add an example sentence to a meaning (contributors)

     Adds the sentence and its English to an existing meaning as `sourced` and queues it
    for blind review (LEX-5, CON-1).

    Args:
        id (UUID):
        body (ExampleContributionInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, ExampleContribution]
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
    body: ExampleContributionInput,
) -> Response[Union[Error, ExampleContribution]]:
    """Add an example sentence to a meaning (contributors)

     Adds the sentence and its English to an existing meaning as `sourced` and queues it
    for blind review (LEX-5, CON-1).

    Args:
        id (UUID):
        body (ExampleContributionInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, ExampleContribution]]
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
    body: ExampleContributionInput,
) -> Optional[Union[Error, ExampleContribution]]:
    """Add an example sentence to a meaning (contributors)

     Adds the sentence and its English to an existing meaning as `sourced` and queues it
    for blind review (LEX-5, CON-1).

    Args:
        id (UUID):
        body (ExampleContributionInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, ExampleContribution]
    """

    return (
        await asyncio_detailed(
            id=id,
            client=client,
            body=body,
        )
    ).parsed
