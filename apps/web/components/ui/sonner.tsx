"use client"

import { Toaster as Sonner, type ToasterProps } from "sonner"

/** shadcn's Toaster without next-themes: follows the system colour scheme. */
export function Toaster(props: ToasterProps) {
  return (
    <Sonner
      theme="system"
      className="toaster group"
      style={
        {
          "--normal-bg": "var(--popover)",
          "--normal-text": "var(--popover-foreground)",
          "--normal-border": "var(--border)",
        } as React.CSSProperties
      }
      {...props}
    />
  )
}
