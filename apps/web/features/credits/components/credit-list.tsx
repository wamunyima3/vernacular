import Link from "next/link"

import { ErrorCard } from "@/components/error-card"
import { NoticeCard } from "@/components/notice-card"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

import { getCredits } from "../get-credits"

/** Everyone credited for accepted work in a language, most credited first (CON-5, US-19). */
export async function CreditList({ language }: { language: string }) {
  const result = await getCredits(language)
  if (!result.ok) return <ErrorCard title="Credits could not be loaded" message={result.message} />
  if (result.value.length === 0) {
    return (
      <NoticeCard
        title="No one to credit yet."
        actions={
          <Button nativeButton={false} render={<Link href="/contribute" />}>
            Contribute a word
          </Button>
        }
      >
        Credit appears here once a contributed word is verified or a review settles an item, for
        volunteers who choose it on their account page.
      </NoticeCard>
    )
  }
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead className="text-right">Words verified</TableHead>
          <TableHead className="text-right">Reviews</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {result.value.map((c, i) => (
          <TableRow key={`${c.name}-${i}`}>
            <TableCell className="max-w-56 truncate font-medium">{c.name}</TableCell>
            <TableCell className="text-right tabular-nums">{c.words}</TableCell>
            <TableCell className="text-right tabular-nums">{c.reviews}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}

export function CreditListSkeleton() {
  return <Skeleton className="h-48 w-full" />
}
