import type { Schemas } from "@vernacular/api-client"
import Link from "next/link"

import { Badge } from "@/components/ui/badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

/** The register as a table. Rows link to the edit page. */
export function SourceTable({ sources }: { sources: Schemas["Source"][] }) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Kind</TableHead>
          <TableHead>Licence</TableHead>
          <TableHead className="hidden sm:table-cell">Rights holder</TableHead>
          <TableHead>Permission</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {sources.map((s) => (
          <TableRow key={s.id}>
            <TableCell className="max-w-48 truncate font-medium">
              <Link href={`/sources/${s.id}`} className="underline-offset-4 hover:underline">
                {s.name}
              </Link>
            </TableCell>
            <TableCell className="capitalize">{s.kind}</TableCell>
            <TableCell>{s.licence}</TableCell>
            <TableCell className="hidden max-w-48 truncate sm:table-cell">
              {s.rights_holder}
            </TableCell>
            <TableCell>
              <Badge variant={s.has_permission_document ? "default" : "outline"}>
                {s.has_permission_document ? "On file" : "Missing"}
              </Badge>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}
