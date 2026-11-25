import type { Metadata } from "next"
import { Atkinson_Hyperlegible_Mono, Atkinson_Hyperlegible_Next } from "next/font/google"

import "./globals.css"

// Atkinson Hyperlegible keeps look-alike letters apart (I/l/1, rn/m): words in unfamiliar
// languages must be read exactly (PRODUCT.md). latin-ext covers names and diacritics.
const sans = Atkinson_Hyperlegible_Next({
  variable: "--font-sans",
  subsets: ["latin", "latin-ext"],
})
const mono = Atkinson_Hyperlegible_Mono({
  variable: "--font-mono-face",
  subsets: ["latin", "latin-ext"],
})

export const metadata: Metadata = {
  title: { default: "Vernacular", template: "%s · Vernacular" },
  description:
    "Verified words, phrases and translations for Zambian languages, through an open API.",
}

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="en" className={`${sans.variable} ${mono.variable} h-full antialiased`}>
      <body className="flex min-h-full flex-col">{children}</body>
    </html>
  )
}
