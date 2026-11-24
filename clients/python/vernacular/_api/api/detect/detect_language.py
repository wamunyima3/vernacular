from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response

from ...models.detect_language_body import DetectLanguageBody
from ...models.detection import Detection
from ...models.error import Error


def _get_kwargs(
    *,
    body: DetectLanguageBody,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/v1/detect",
    }

    _kwargs["json"] = body.to_dict()

    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Detection, Error]:
    if response.status_code == 200:
        response_200 = Detection.from_dict(response.json())

        return response_200

    if response.status_code == 400:
        response_400 = Error.from_dict(response.json())

        return response_400

    if response.status_code == 401:
        response_401 = Error.from_dict(response.json())

        return response_401

    if response.status_code == 429:
        response_429 = Error.from_dict(response.json())

        return response_429

    response_default = Error.from_dict(response.json())

    return response_default


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[Detection, Error]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    body: DetectLanguageBody,
) -> Response[Union[Detection, Error]]:
    """Detect the language of a text

     Ranks the supported languages, and English, by how likely the text is written in
    each, with confidences between 0 and 1 that sum to 1 (TRN-1, US-09). Short texts
    and closely related languages give close scores, so check the margin. Empty when
    the text has no letters. Languages without enough data yet are not detected
    (ADR-0010).

    Args:
        body (DetectLanguageBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Detection, Error]]
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
    body: DetectLanguageBody,
) -> Optional[Union[Detection, Error]]:
    """Detect the language of a text

     Ranks the supported languages, and English, by how likely the text is written in
    each, with confidences between 0 and 1 that sum to 1 (TRN-1, US-09). Short texts
    and closely related languages give close scores, so check the margin. Empty when
    the text has no letters. Languages without enough data yet are not detected
    (ADR-0010).

    Args:
        body (DetectLanguageBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Detection, Error]
    """

    return sync_detailed(
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    body: DetectLanguageBody,
) -> Response[Union[Detection, Error]]:
    """Detect the language of a text

     Ranks the supported languages, and English, by how likely the text is written in
    each, with confidences between 0 and 1 that sum to 1 (TRN-1, US-09). Short texts
    and closely related languages give close scores, so check the margin. Empty when
    the text has no letters. Languages without enough data yet are not detected
    (ADR-0010).

    Args:
        body (DetectLanguageBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Detection, Error]]
    """

    kwargs = _get_kwargs(
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    body: DetectLanguageBody,
) -> Optional[Union[Detection, Error]]:
    """Detect the language of a text

     Ranks the supported languages, and English, by how likely the text is written in
    each, with confidences between 0 and 1 that sum to 1 (TRN-1, US-09). Short texts
    and closely related languages give close scores, so check the margin. Empty when
    the text has no letters. Languages without enough data yet are not detected
    (ADR-0010).

    Args:
        body (DetectLanguageBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Detection, Error]
    """

    return (
        await asyncio_detailed(
            client=client,
            body=body,
        )
    ).parsed
