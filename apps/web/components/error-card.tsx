import { NoticeCard } from "./notice-card"

/** A failed load, shown where the content would be, with one next step. */
export function ErrorCard({ title, message }: { title: string; message: string }) {
  return (
    <NoticeCard title={title} role="alert">
      {message} Reload the page to try again.
    </NoticeCard>
  )
}
