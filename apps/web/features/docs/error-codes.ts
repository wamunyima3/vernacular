import type { Schemas } from "@vernacular/api-client"

type ErrorCode = Schemas["ErrorCode"]

/**
 * What each error code means and what to do about it. Typed against the contract,
 * so adding a code to `contracts/openapi.yaml` fails the type check until it is described here.
 */
export const errorCodes: Record<ErrorCode, { status: number; meaning: string; action: string }> = {
  bad_request: {
    status: 400,
    meaning: "A parameter is missing, malformed or out of range.",
    action: "Read the message, fix the request and send it again.",
  },
  unauthorized: {
    status: 401,
    meaning: "The API key is missing, unknown or revoked.",
    action: "Send a valid key in the Authorization header.",
  },
  forbidden: {
    status: 403,
    meaning: "The key or account is not allowed to do this.",
    action: "Use an account with the right role.",
  },
  not_found: {
    status: 404,
    meaning: "The route or item does not exist, or no verified item matches.",
    action: "Check the ID or path. To see unverified items, set min_status.",
  },
  method_not_allowed: {
    status: 405,
    meaning: "The route exists but does not accept this HTTP method.",
    action: "Use the method listed in the API reference.",
  },
  rate_limited: {
    status: 429,
    meaning: "The key or IP address has used up its request allowance.",
    action: "Wait until the time in the RateLimit-Reset header.",
  },
  internal: {
    status: 500,
    meaning: "Something failed on our side. Details are logged, not returned.",
    action: "Retry later. If it keeps happening, report it with the X-Request-Id header.",
  },
  unavailable: {
    status: 503,
    meaning: "A dependency such as the database is unreachable.",
    action: "Retry with backoff.",
  },
}
