from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.error import Error
from ...models.recheck_batch_list import RecheckBatchList
from ...types import Unset


def _get_kwargs(
    *,
    quarter: Union[Unset, str] = UNSET,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["quarter"] = quarter

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/v1/rechecks/batches",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, RecheckBatchList]:
    if response.status_code == 200:
        response_200 = RecheckBatchList.from_dict(response.json())

        return response_200

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Error, RecheckBatchList]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    quarter: Union[Unset, str] = UNSET,
) -> Response[Union[Error, RecheckBatchList]]:
    """A quarter's re-checks per batch (admins)

     Each batch, the source and model its items came from in a language, with how many
    of its items were drawn, checked and failed, worst first. A batch failing more than
    2% of its checks can be sent back to review (DQ-6).

    Args:
        quarter (Union[Unset, str]): A calendar quarter, UTC Example: 2026-Q4.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, RecheckBatchList]]
    """

    kwargs = _get_kwargs(
        quarter=quarter,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: AuthenticatedClient,
    quarter: Union[Unset, str] = UNSET,
) -> Optional[Union[Error, RecheckBatchList]]:
    """A quarter's re-checks per batch (admins)

     Each batch, the source and model its items came from in a language, with how many
    of its items were drawn, checked and failed, worst first. A batch failing more than
    2% of its checks can be sent back to review (DQ-6).

    Args:
        quarter (Union[Unset, str]): A calendar quarter, UTC Example: 2026-Q4.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, RecheckBatchList]
    """

    return sync_detailed(
        client=client,
        quarter=quarter,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    quarter: Union[Unset, str] = UNSET,
) -> Response[Union[Error, RecheckBatchList]]:
    """A quarter's re-checks per batch (admins)

     Each batch, the source and model its items came from in a language, with how many
    of its items were drawn, checked and failed, worst first. A batch failing more than
    2% of its checks can be sent back to review (DQ-6).

    Args:
        quarter (Union[Unset, str]): A calendar quarter, UTC Example: 2026-Q4.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, RecheckBatchList]]
    """

    kwargs = _get_kwargs(
        quarter=quarter,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    quarter: Union[Unset, str] = UNSET,
) -> Optional[Union[Error, RecheckBatchList]]:
    """A quarter's re-checks per batch (admins)

     Each batch, the source and model its items came from in a language, with how many
    of its items were drawn, checked and failed, worst first. A batch failing more than
    2% of its checks can be sent back to review (DQ-6).

    Args:
        quarter (Union[Unset, str]): A calendar quarter, UTC Example: 2026-Q4.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, RecheckBatchList]
    """

    return (
        await asyncio_detailed(
            client=client,
            quarter=quarter,
        )
    ).parsed
