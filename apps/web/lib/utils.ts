export { cn } from "cn"

/** "licensed" → "Licensed", for option labels (CSS `capitalize` does not reach native options). */
export const capitalise = (s: string) => (s[0]?.toUpperCase() ?? "") + s.slice(1)
