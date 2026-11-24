# Translation service contract

The machine translation model runs as a separate service (P3-03), so the API stays a
small Go binary and the model can move to a GPU host or change without an API release.
The API calls it only for `POST /v1/translate` with `machine: true` when no verified
phrase matches (TRN-3). Set `TRANSLATION_SERVICE_URL` to its base URL; empty means none.

## `POST {base}/translate`

Request, JSON:

| Field  | Type   |                                                              |
| ------ | ------ | ------------------------------------------------------------ |
| `text` | string | Up to 500 characters                                         |
| `from` | string | Language code, `eng` or a Zambian language (`bem`, `nya`, …) |
| `to`   | string | Language code; one side is always `eng`                      |

Response:

- `200` with JSON `{ "translation": string, "confidence": number, "model": string }`.
  `confidence` is between 0 and 1. `model` names the model and its version, as it will
  appear to API users (for example `vernacular-mt-0.1`).
- `422` when the service does not translate that language pair. The API leaves the
  machine translation out without logging an error.
- Anything else is an error: the API logs it and leaves the machine translation out.
  Verified results never depend on the service.

The API waits up to 5 seconds, since the caller is waiting too. The service should never return a verified phrase as its
own output, and must not be trained on the held-out test sets (DQ-10, `ml/README.md`).
