const dateFormat = new Intl.DateTimeFormat("en-GB", { dateStyle: "medium" })
const countFormat = new Intl.NumberFormat("en-GB")

/** A date as "24 Nov 2026", the one date style the app uses. */
export const formatDate = (iso: string) => dateFormat.format(new Date(iso))

/** part as a whole percentage of whole, 0 when whole is 0. */
export const percentOf = (part: number, whole: number) =>
  whole === 0 ? 0 : Math.round((100 * part) / whole)

/** A count as "1,204". */
export const formatCount = (n: number) => countFormat.format(n)
