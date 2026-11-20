import { notFound } from "next/navigation"

import { Separator } from "@/components/ui/separator"
import { api } from "@/lib/api"

import { sessionWithAdmin } from "../../auth/admin"
import { deleteSource, updateSource, uploadPermission } from "../actions"
import { AdminOnly } from "./admin-only"
import { DeleteSource } from "./delete-source"
import { PermissionUpload } from "./permission-upload"
import { SourceForm } from "./source-form"

export async function SourceEditor({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params
  const { isAdmin } = await sessionWithAdmin(`/sources/${id}`)
  if (!isAdmin) return <AdminOnly />

  const { data: source, response } = await api.GET("/v1/sources/{id}", {
    params: { path: { id } },
  })
  if (!source) {
    if (response.status === 404 || response.status === 400) notFound()
    throw new Error(`Loading source ${id} failed with ${response.status}`)
  }

  return (
    <div className="flex flex-col gap-8">
      <h1 className="text-2xl font-semibold text-balance">{source.name}</h1>
      <SourceForm action={updateSource.bind(null, id)} source={source} submitLabel="Save changes" />
      <Separator />
      <PermissionUpload
        action={uploadPermission.bind(null, id)}
        hasDocument={source.has_permission_document}
        downloadHref={`/sources/${id}/permission`}
      />
      <Separator />
      <DeleteSource name={source.name} action={deleteSource.bind(null, id)} />
    </div>
  )
}
