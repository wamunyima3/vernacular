import { Skeleton } from "@/components/ui/skeleton"

import { LanguageGrid } from "./language-grid"

/** Structural placeholder matching LanguageList's layout. */
export function LanguageListSkeleton() {
  return (
    <div className="flex flex-col gap-3" aria-busy="true" aria-label="Loading languages">
      <Skeleton className="h-4 w-16" />
      <LanguageGrid>
        {Array.from({ length: 4 }, (_, i) => (
          <Skeleton key={i} className="h-24 rounded-xl" />
        ))}
      </LanguageGrid>
    </div>
  )
}
