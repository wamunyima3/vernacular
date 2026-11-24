from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.error import Error
from ...models.lookup_results import LookupResults
from ...models.min_status import MinStatus
from ...types import Unset


def _get_kwargs(
    *,
    q: str,
    from_: str,
    to: Union[Unset, list[str]] = UNSET,
    limit: Union[Unset, int] = 5,
    min_status: Union[Unset, MinStatus] = UNSET,
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["q"] = q

    params["from"] = from_

    json_to: Union[Unset, list[str]] = UNSET
    if not isinstance(to, Unset):
        json_to = to

    params["to"] = json_to

    params["limit"] = limit

    json_min_status: Union[Unset, str] = UNSET
    if not isinstance(min_status, Unset):
        json_min_status = min_status.value

    params["min_status"] = json_min_status

    params["variety"] = variety

    params["release"] = release

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/v1/lookup",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, LookupResults]:
    if response.status_code == 200:
        response_200 = LookupResults.from_dict(response.json())

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
) -> Response[Union[Error, LookupResults]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    q: str,
    from_: str,
    to: Union[Unset, list[str]] = UNSET,
    limit: Union[Unset, int] = 5,
    min_status: Union[Unset, MinStatus] = UNSET,
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Response[Union[Error, LookupResults]]:
    """Translate a word between English and Zambian languages

     From English (`from=eng`), finds concepts whose English gloss matches the query and
    returns their entries in the `to` languages (LEX-2, US-04). From a Zambian language,
    searches its entries as `/v1/entries/search` does and returns the English meaning of
    each, grouped by concept. Each result says what matched; best matches first.

    Args:
        q (str):
        from_ (str): ISO 639-3 code Example: bem.
        to (Union[Unset, list[str]]):
        limit (Union[Unset, int]):  Default: 5.
        min_status (Union[Unset, MinStatus]):
        variety (Union[Unset, str]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, LookupResults]]
    """

    kwargs = _get_kwargs(
        q=q,
        from_=from_,
        to=to,
        limit=limit,
        min_status=min_status,
        variety=variety,
        release=release,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: AuthenticatedClient,
    q: str,
    from_: str,
    to: Union[Unset, list[str]] = UNSET,
    limit: Union[Unset, int] = 5,
    min_status: Union[Unset, MinStatus] = UNSET,
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Optional[Union[Error, LookupResults]]:
    """Translate a word between English and Zambian languages

     From English (`from=eng`), finds concepts whose English gloss matches the query and
    returns their entries in the `to` languages (LEX-2, US-04). From a Zambian language,
    searches its entries as `/v1/entries/search` does and returns the English meaning of
    each, grouped by concept. Each result says what matched; best matches first.

    Args:
        q (str):
        from_ (str): ISO 639-3 code Example: bem.
        to (Union[Unset, list[str]]):
        limit (Union[Unset, int]):  Default: 5.
        min_status (Union[Unset, MinStatus]):
        variety (Union[Unset, str]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, LookupResults]
    """

    return sync_detailed(
        client=client,
        q=q,
        from_=from_,
        to=to,
        limit=limit,
        min_status=min_status,
        variety=variety,
        release=release,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    q: str,
    from_: str,
    to: Union[Unset, list[str]] = UNSET,
    limit: Union[Unset, int] = 5,
    min_status: Union[Unset, MinStatus] = UNSET,
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Response[Union[Error, LookupResults]]:
    """Translate a word between English and Zambian languages

     From English (`from=eng`), finds concepts whose English gloss matches the query and
    returns their entries in the `to` languages (LEX-2, US-04). From a Zambian language,
    searches its entries as `/v1/entries/search` does and returns the English meaning of
    each, grouped by concept. Each result says what matched; best matches first.

    Args:
        q (str):
        from_ (str): ISO 639-3 code Example: bem.
        to (Union[Unset, list[str]]):
        limit (Union[Unset, int]):  Default: 5.
        min_status (Union[Unset, MinStatus]):
        variety (Union[Unset, str]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, LookupResults]]
    """

    kwargs = _get_kwargs(
        q=q,
        from_=from_,
        to=to,
        limit=limit,
        min_status=min_status,
        variety=variety,
        release=release,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    q: str,
    from_: str,
    to: Union[Unset, list[str]] = UNSET,
    limit: Union[Unset, int] = 5,
    min_status: Union[Unset, MinStatus] = UNSET,
    variety: Union[Unset, str] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Optional[Union[Error, LookupResults]]:
    """Translate a word between English and Zambian languages

     From English (`from=eng`), finds concepts whose English gloss matches the query and
    returns their entries in the `to` languages (LEX-2, US-04). From a Zambian language,
    searches its entries as `/v1/entries/search` does and returns the English meaning of
    each, grouped by concept. Each result says what matched; best matches first.

    Args:
        q (str):
        from_ (str): ISO 639-3 code Example: bem.
        to (Union[Unset, list[str]]):
        limit (Union[Unset, int]):  Default: 5.
        min_status (Union[Unset, MinStatus]):
        variety (Union[Unset, str]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, LookupResults]
    """

    return (
        await asyncio_detailed(
            client=client,
            q=q,
            from_=from_,
            to=to,
            limit=limit,
            min_status=min_status,
            variety=variety,
            release=release,
        )
    ).parsed
