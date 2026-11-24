const dateFormat = new Intl.DateTimeFormat("en-GB", { dateStyle: "medium" })

/** A date as "24 Nov 2026", the one date style the app uses. */
export const formatDate = (iso: string) => dateFormat.format(new Date(iso))
