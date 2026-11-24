import { ConfirmAction } from "@/components/confirm-action"

import { apiForVisitor } from "../../auth/session"
import { deleteClip } from "../actions"
import { Recorder } from "./recorder"

/**
 * An entry's pronunciations (LEX-8, US-33), and a recorder for a speaker whose consent
 * covers the entry's language. Shows nothing while the server has audio switched off.
 */
export async function ClipList({ entryId }: { entryId: string }) {
  const api = await apiForVisitor()
  const { data } = await api.GET("/v1/audio", { params: { query: { entry_id: entryId } } })
  if (!data) return null
  const { data: clips, can_record: canRecord } = data
  if (clips.length === 0 && !canRecord) return null

  return (
    <section className="flex flex-col gap-3" aria-labelledby="pronunciation">
      <h2 id="pronunciation" className="text-lg font-semibold text-balance">
        Pronunciation
      </h2>
      {clips.length > 0 && (
        <ul className="flex flex-col gap-2">
          {clips.map((c, i) => (
            <li key={c.id} className="flex flex-wrap items-center gap-3">
              <audio
                controls
                preload="none"
                src={`/audio/${c.id}`}
                aria-label={`Pronunciation ${i + 1}`}
                className="h-9 max-w-full"
              />
              {c.mine && (
                <ConfirmAction
                  label="Delete"
                  title="Delete your recording?"
                  description="The recording is removed for everyone. You can record the word again."
                  action={deleteClip.bind(null, entryId, c.id)}
                  destructive
                />
              )}
            </li>
          ))}
        </ul>
      )}
      {canRecord && <Recorder entryId={entryId} />}
    </section>
  )
}
