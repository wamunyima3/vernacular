/** A fetch result the UI can show in place instead of throwing. */
export type Settled<T> = { ok: true; value: T } | { ok: false; message: string }

/** Runs fn, turning a thrown error into a message. */
export async function settle<T>(fn: () => Promise<T>): Promise<Settled<T>> {
  try {
    return { ok: true, value: await fn() }
  } catch (e) {
    return {
      ok: false,
      message: e instanceof Error && e.message ? e.message : "The API could not be reached.",
    }
  }
}
