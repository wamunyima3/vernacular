from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.error import Error
from ...models.translate_body import TranslateBody
from ...models.translation_results import TranslationResults
from ...types import Unset


def _get_kwargs(
    *,
    body: TranslateBody,
    release: Union[Unset, str] = UNSET,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    params: dict[str, Any] = {}

    params["release"] = release

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/v1/translate",
        "params": params,
    }

    _kwargs["json"] = body.to_dict()

    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, TranslationResults]:
    if response.status_code == 200:
        response_200 = TranslationResults.from_dict(response.json())

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
) -> Response[Union[Error, TranslationResults]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    body: TranslateBody,
    release: Union[Unset, str] = UNSET,
) -> Response[Union[Error, TranslationResults]]:
    """Translate a short text between English and a Zambian language

     Exact match, ignoring case and accents, against verified phrases (TRN-2, US-10).
    Each match has `method: verified_phrase`. No match returns an empty list, never a
    guess in `data`. With `machine: true` and no verified match, `machine` may hold a
    machine translation, marked `status: machine` with a confidence score (TRN-3,
    US-31). It is left out when the server has no translation service or the service
    cannot translate the pair.

    Args:
        release (Union[Unset, str]):  Example: 2027.1.
        body (TranslateBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, TranslationResults]]
    """

    kwargs = _get_kwargs(
        body=body,
        release=release,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: AuthenticatedClient,
    body: TranslateBody,
    release: Union[Unset, str] = UNSET,
) -> Optional[Union[Error, TranslationResults]]:
    """Translate a short text between English and a Zambian language

     Exact match, ignoring case and accents, against verified phrases (TRN-2, US-10).
    Each match has `method: verified_phrase`. No match returns an empty list, never a
    guess in `data`. With `machine: true` and no verified match, `machine` may hold a
    machine translation, marked `status: machine` with a confidence score (TRN-3,
    US-31). It is left out when the server has no translation service or the service
    cannot translate the pair.

    Args:
        release (Union[Unset, str]):  Example: 2027.1.
        body (TranslateBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, TranslationResults]
    """

    return sync_detailed(
        client=client,
        body=body,
        release=release,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    body: TranslateBody,
    release: Union[Unset, str] = UNSET,
) -> Response[Union[Error, TranslationResults]]:
    """Translate a short text between English and a Zambian language

     Exact match, ignoring case and accents, against verified phrases (TRN-2, US-10).
    Each match has `method: verified_phrase`. No match returns an empty list, never a
    guess in `data`. With `machine: true` and no verified match, `machine` may hold a
    machine translation, marked `status: machine` with a confidence score (TRN-3,
    US-31). It is left out when the server has no translation service or the service
    cannot translate the pair.

    Args:
        release (Union[Unset, str]):  Example: 2027.1.
        body (TranslateBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, TranslationResults]]
    """

    kwargs = _get_kwargs(
        body=body,
        release=release,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    body: TranslateBody,
    release: Union[Unset, str] = UNSET,
) -> Optional[Union[Error, TranslationResults]]:
    """Translate a short text between English and a Zambian language

     Exact match, ignoring case and accents, against verified phrases (TRN-2, US-10).
    Each match has `method: verified_phrase`. No match returns an empty list, never a
    guess in `data`. With `machine: true` and no verified match, `machine` may hold a
    machine translation, marked `status: machine` with a confidence score (TRN-3,
    US-31). It is left out when the server has no translation service or the service
    cannot translate the pair.

    Args:
        release (Union[Unset, str]):  Example: 2027.1.
        body (TranslateBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, TranslationResults]
    """

    return (
        await asyncio_detailed(
            client=client,
            body=body,
            release=release,
        )
    ).parsed
