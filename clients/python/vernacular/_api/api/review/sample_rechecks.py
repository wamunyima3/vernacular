from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response

from ...models.error import Error
from ...models.recheck_sample import RecheckSample
from ...models.recheck_sample_input import RecheckSampleInput


def _get_kwargs(
    *,
    body: RecheckSampleInput,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/v1/rechecks",
    }

    _kwargs["json"] = body.to_dict()

    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, RecheckSample]:
    if response.status_code == 200:
        response_200 = RecheckSample.from_dict(response.json())

        return response_200

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Error, RecheckSample]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    body: RecheckSampleInput,
) -> Response[Union[Error, RecheckSample]]:
    """Draw a quarter's re-check sample (admins)

     Draws 5% of each language's verified items, at least one, to be checked again by
    senior reviewers (DQ-6). A quarter is drawn once, so calling again returns
    `sampled: 0`. The worker draws each quarter on its own; this starts it by hand.

    Args:
        body (RecheckSampleInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, RecheckSample]]
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
    body: RecheckSampleInput,
) -> Optional[Union[Error, RecheckSample]]:
    """Draw a quarter's re-check sample (admins)

     Draws 5% of each language's verified items, at least one, to be checked again by
    senior reviewers (DQ-6). A quarter is drawn once, so calling again returns
    `sampled: 0`. The worker draws each quarter on its own; this starts it by hand.

    Args:
        body (RecheckSampleInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, RecheckSample]
    """

    return sync_detailed(
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    body: RecheckSampleInput,
) -> Response[Union[Error, RecheckSample]]:
    """Draw a quarter's re-check sample (admins)

     Draws 5% of each language's verified items, at least one, to be checked again by
    senior reviewers (DQ-6). A quarter is drawn once, so calling again returns
    `sampled: 0`. The worker draws each quarter on its own; this starts it by hand.

    Args:
        body (RecheckSampleInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, RecheckSample]]
    """

    kwargs = _get_kwargs(
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    body: RecheckSampleInput,
) -> Optional[Union[Error, RecheckSample]]:
    """Draw a quarter's re-check sample (admins)

     Draws 5% of each language's verified items, at least one, to be checked again by
    senior reviewers (DQ-6). A quarter is drawn once, so calling again returns
    `sampled: 0`. The worker draws each quarter on its own; this starts it by hand.

    Args:
        body (RecheckSampleInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, RecheckSample]
    """

    return (
        await asyncio_detailed(
            client=client,
            body=body,
        )
    ).parsed
