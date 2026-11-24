"use client"

import { useState } from "react"

import { Button } from "@/components/ui/button"

/** Copies text to the clipboard and says so until the text changes. */
export function CopyButton({
  text,
  label,
  size = "default",
}: {
  text: string
  label?: string // the accessible name, when "Copy" alone would be unclear
  size?: "default" | "sm"
}) {
  const [copied, setCopied] = useState<string>()
  return (
    <Button
      type="button"
      variant="outline"
      size={size}
      aria-label={label}
      onClick={() => navigator.clipboard.writeText(text).then(() => setCopied(text))}
    >
      {copied === text ? "Copied" : "Copy"}
    </Button>
  )
}
