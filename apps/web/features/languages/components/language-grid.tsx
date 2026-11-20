/** Responsive grid shared by the language list and its loading skeleton. */
export function LanguageGrid({ children }: { children: React.ReactNode }) {
  return <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">{children}</div>
}
