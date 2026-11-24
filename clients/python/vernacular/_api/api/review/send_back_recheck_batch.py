from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response

from ...models.error import Error
from ...models.recheck_batch_input import RecheckBatchInput
from ...models.recheck_batch_sent_back import RecheckBatchSentBack


def _get_kwargs(
    *,
    body: RecheckBatchInput,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/v1/rechecks/batches/send-back",
    }

    _kwargs["json"] = body.to_dict()

    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, RecheckBatchSentBack]:
    if response.status_code == 200:
        response_200 = RecheckBatchSentBack.from_dict(response.json())

        return response_200

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Error, RecheckBatchSentBack]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    body: RecheckBatchInput,
) -> Response[Union[Error, RecheckBatchSentBack]]:
    """Send a failing batch back to review (admins)

     Every verified item of the batch goes back to the review queue, where only reviews
    made after this count. Refused (`conflict`) unless more than 2% of the batch's
    checks that quarter failed, or when it was already sent back.

    Args:
        body (RecheckBatchInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, RecheckBatchSentBack]]
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
    body: RecheckBatchInput,
) -> Optional[Union[Error, RecheckBatchSentBack]]:
    """Send a failing batch back to review (admins)

     Every verified item of the batch goes back to the review queue, where only reviews
    made after this count. Refused (`conflict`) unless more than 2% of the batch's
    checks that quarter failed, or when it was already sent back.

    Args:
        body (RecheckBatchInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, RecheckBatchSentBack]
    """

    return sync_detailed(
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    body: RecheckBatchInput,
) -> Response[Union[Error, RecheckBatchSentBack]]:
    """Send a failing batch back to review (admins)

     Every verified item of the batch goes back to the review queue, where only reviews
    made after this count. Refused (`conflict`) unless more than 2% of the batch's
    checks that quarter failed, or when it was already sent back.

    Args:
        body (RecheckBatchInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, RecheckBatchSentBack]]
    """

    kwargs = _get_kwargs(
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    body: RecheckBatchInput,
) -> Optional[Union[Error, RecheckBatchSentBack]]:
    """Send a failing batch back to review (admins)

     Every verified item of the batch goes back to the review queue, where only reviews
    made after this count. Refused (`conflict`) unless more than 2% of the batch's
    checks that quarter failed, or when it was already sent back.

    Args:
        body (RecheckBatchInput):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, RecheckBatchSentBack]
    """

    return (
        await asyncio_detailed(
            client=client,
            body=body,
        )
    ).parsed
