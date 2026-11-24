import Link from "next/link"

import { NoticeCard } from "@/components/notice-card"
import { Button } from "@/components/ui/button"

/** Shown when a public page needs the live API and the visitor has no access to it. */
export function SignInNotice({
  title,
  next,
  children,
}: {
  title: string
  next: string
  children: React.ReactNode
}) {
  return (
    <NoticeCard
      title={title}
      actions={
        <Button
          nativeButton={false}
          render={<Link href={`/sign-in?${new URLSearchParams({ next })}`} />}
        >
          Sign in
        </Button>
      }
    >
      {children}
    </NoticeCard>
  )
}
