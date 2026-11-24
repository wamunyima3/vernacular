const dateFormat = new Intl.DateTimeFormat("en-GB", { dateStyle: "medium" })
const countFormat = new Intl.NumberFormat("en-GB")

/** A date as "24 Nov 2026", the one date style the app uses. */
export const formatDate = (iso: string) => dateFormat.format(new Date(iso))

/** A count as "1,204". */
export const formatCount = (n: number) => countFormat.format(n)
