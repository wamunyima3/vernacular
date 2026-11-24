from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ...client import AuthenticatedClient, Client
from ...types import Response, UNSET

from ...models.error import Error
from ...models.export_phrase_pack_format import ExportPhrasePackFormat
from ...models.export_phrase_pack_response_200 import ExportPhrasePackResponse200
from ...models.min_status import MinStatus
from ...types import Unset


def _get_kwargs(
    *,
    pack: str,
    lang: str,
    format_: ExportPhrasePackFormat,
    min_status: Union[Unset, MinStatus] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["pack"] = pack

    params["lang"] = lang

    json_format_ = format_.value
    params["format"] = json_format_

    json_min_status: Union[Unset, str] = UNSET
    if not isinstance(min_status, Unset):
        json_min_status = min_status.value

    params["min_status"] = json_min_status

    params["release"] = release

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/v1/phrases/export",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Union[Error, ExportPhrasePackResponse200]:
    if response.status_code == 200:
        response_200 = ExportPhrasePackResponse200.from_dict(response.json())

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
) -> Response[Union[Error, ExportPhrasePackResponse200]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: AuthenticatedClient,
    pack: str,
    lang: str,
    format_: ExportPhrasePackFormat,
    min_status: Union[Unset, MinStatus] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Response[Union[Error, ExportPhrasePackResponse200]]:
    """A phrase pack as a localisation file

     The pack as a file an i18n pipeline reads (PHR-3, US-12), keyed by the English
    source: `json` is a flat object from English to translation, `po` a gettext
    catalogue, and `xliff` an XLIFF 2.0 document whose segment state follows the trust
    status (`verified` is `final`). When an English line has several translations,
    the first is kept. Notes become translator comments.

    Args:
        pack (str):
        lang (str): ISO 639-3 code Example: bem.
        format_ (ExportPhrasePackFormat):
        min_status (Union[Unset, MinStatus]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, ExportPhrasePackResponse200]]
    """

    kwargs = _get_kwargs(
        pack=pack,
        lang=lang,
        format_=format_,
        min_status=min_status,
        release=release,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: AuthenticatedClient,
    pack: str,
    lang: str,
    format_: ExportPhrasePackFormat,
    min_status: Union[Unset, MinStatus] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Optional[Union[Error, ExportPhrasePackResponse200]]:
    """A phrase pack as a localisation file

     The pack as a file an i18n pipeline reads (PHR-3, US-12), keyed by the English
    source: `json` is a flat object from English to translation, `po` a gettext
    catalogue, and `xliff` an XLIFF 2.0 document whose segment state follows the trust
    status (`verified` is `final`). When an English line has several translations,
    the first is kept. Notes become translator comments.

    Args:
        pack (str):
        lang (str): ISO 639-3 code Example: bem.
        format_ (ExportPhrasePackFormat):
        min_status (Union[Unset, MinStatus]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, ExportPhrasePackResponse200]
    """

    return sync_detailed(
        client=client,
        pack=pack,
        lang=lang,
        format_=format_,
        min_status=min_status,
        release=release,
    ).parsed


async def asyncio_detailed(
    *,
    client: AuthenticatedClient,
    pack: str,
    lang: str,
    format_: ExportPhrasePackFormat,
    min_status: Union[Unset, MinStatus] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Response[Union[Error, ExportPhrasePackResponse200]]:
    """A phrase pack as a localisation file

     The pack as a file an i18n pipeline reads (PHR-3, US-12), keyed by the English
    source: `json` is a flat object from English to translation, `po` a gettext
    catalogue, and `xliff` an XLIFF 2.0 document whose segment state follows the trust
    status (`verified` is `final`). When an English line has several translations,
    the first is kept. Notes become translator comments.

    Args:
        pack (str):
        lang (str): ISO 639-3 code Example: bem.
        format_ (ExportPhrasePackFormat):
        min_status (Union[Unset, MinStatus]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[Error, ExportPhrasePackResponse200]]
    """

    kwargs = _get_kwargs(
        pack=pack,
        lang=lang,
        format_=format_,
        min_status=min_status,
        release=release,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: AuthenticatedClient,
    pack: str,
    lang: str,
    format_: ExportPhrasePackFormat,
    min_status: Union[Unset, MinStatus] = UNSET,
    release: Union[Unset, str] = UNSET,
) -> Optional[Union[Error, ExportPhrasePackResponse200]]:
    """A phrase pack as a localisation file

     The pack as a file an i18n pipeline reads (PHR-3, US-12), keyed by the English
    source: `json` is a flat object from English to translation, `po` a gettext
    catalogue, and `xliff` an XLIFF 2.0 document whose segment state follows the trust
    status (`verified` is `final`). When an English line has several translations,
    the first is kept. Notes become translator comments.

    Args:
        pack (str):
        lang (str): ISO 639-3 code Example: bem.
        format_ (ExportPhrasePackFormat):
        min_status (Union[Unset, MinStatus]):
        release (Union[Unset, str]):  Example: 2027.1.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[Error, ExportPhrasePackResponse200]
    """

    return (
        await asyncio_detailed(
            client=client,
            pack=pack,
            lang=lang,
            format_=format_,
            min_status=min_status,
            release=release,
        )
    ).parsed
