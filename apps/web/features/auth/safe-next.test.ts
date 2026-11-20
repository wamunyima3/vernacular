import { describe, expect, it } from "vitest"

import { safeNext } from "./safe-next"

describe("safeNext", () => {
  it.each([
    ["/review", "/review"],
    ["/account?tab=keys", "/account?tab=keys"],
    ["https://evil.example", "/account"],
    ["//evil.example", "/account"],
    ["", "/account"],
    [null, "/account"],
  ])("%s → %s", (input, want) => {
    expect(safeNext(input)).toBe(want)
  })
})
