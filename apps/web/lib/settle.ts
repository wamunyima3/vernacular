import "server-only"

import { connection } from "next/server"

/** A fetch result the UI can show in place instead of throwing. */
export type Settled<T> = { ok: true; value: T } | { ok: false; message: string }

/**
 * Runs fn, a "use cache" function that reads the API, turning a thrown error into a message.
 * A container image is built without the API (apps/web/Dockerfile sets BUILD_WITHOUT_API=1),
 * and a cached function that throws while prerendering fails the build, so there the read
 * waits for a request instead. At run time the variable is unset.
 */
export async function settle<T>(fn: () => Promise<T>): Promise<Settled<T>> {
  if (process.env.BUILD_WITHOUT_API) await connection()
  try {
    return { ok: true, value: await fn() }
  } catch (e) {
    return {
      ok: false,
      message: e instanceof Error && e.message ? e.message : "The API could not be reached.",
    }
  }
}
