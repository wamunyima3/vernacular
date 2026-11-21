import { NoticeCard } from "@/components/notice-card"

/** Shown to signed-in users who are not admins; children say what admins do here. */
export function AdminOnly({ children }: { children: React.ReactNode }) {
  return <NoticeCard title="Admins only">{children}</NoticeCard>
}
