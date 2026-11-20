import { NoticeCard } from "@/components/notice-card"

/** Shown to signed-in users who are not admins. */
export function AdminOnly() {
  return (
    <NoticeCard title="Admins only">
      Only admins manage the source register. Ask an admin if a source is missing.
    </NoticeCard>
  )
}
