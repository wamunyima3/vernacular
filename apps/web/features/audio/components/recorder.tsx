"use client"

import { startTransition, useActionState, useState } from "react"

import { Button } from "@/components/ui/button"
import { FieldError } from "@/components/ui/field"

import { recordClip, type RecordState } from "../actions"

/** A spoken word needs a few seconds; stop well before the 1 MB limit. */
const maxSeconds = 8

type Take = { blob: Blob; url: string }

/**
 * Records the entry in the browser with MediaRecorder, plays the take back, and saves it
 * (LEX-8). Saving replaces the speaker's earlier recording of the entry.
 */
export function Recorder({ entryId }: { entryId: string }) {
  const [state, save, saving] = useActionState<RecordState, FormData>(
    recordClip.bind(null, entryId),
    {},
  )
  const [recorder, setRecorder] = useState<MediaRecorder>() // set while recording
  const [take, setTake] = useState<Take>()
  const [micError, setMicError] = useState<string>()
  const recording = recorder !== undefined

  async function start() {
    setMicError(undefined)
    let stream: MediaStream
    try {
      stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    } catch {
      setMicError("Allow this site to use your microphone, then try again.")
      return
    }
    const r = new MediaRecorder(stream)
    const chunks: Blob[] = []
    r.ondataavailable = (e) => chunks.push(e.data)
    r.onstop = () => {
      stream.getTracks().forEach((t) => t.stop())
      const blob = new Blob(chunks, { type: r.mimeType })
      setTake((old) => {
        if (old) URL.revokeObjectURL(old.url)
        return { blob, url: URL.createObjectURL(blob) }
      })
      setRecorder(undefined)
    }
    r.start()
    setRecorder(r)
    setTimeout(() => r.state === "recording" && r.stop(), maxSeconds * 1000)
  }

  function submit() {
    if (!take) return
    const form = new FormData()
    form.set("clip", take.blob)
    startTransition(() => save(form))
  }

  const error = micError ?? state.error
  return (
    <div className="flex flex-col gap-2">
      <p className="text-sm text-pretty text-muted-foreground">
        Say the word once, clearly. Recording stops after {maxSeconds} seconds.
      </p>
      <div className="flex flex-wrap items-center gap-2">
        {recording ? (
          <Button variant="outline" onClick={() => recorder?.stop()}>
            Stop
          </Button>
        ) : (
          <Button variant="outline" onClick={start} disabled={saving}>
            {take ? "Record again" : "Record"}
          </Button>
        )}
        {take && !recording && (
          <>
            <audio controls src={take.url} aria-label="Your recording" className="h-9 max-w-full" />
            <Button onClick={submit} disabled={saving}>
              {saving ? "Saving…" : "Save recording"}
            </Button>
          </>
        )}
      </div>
      {recording && (
        <p role="status" className="text-sm text-pretty">
          Recording…
        </p>
      )}
      {state.saved && !saving && (
        <p role="status" className="text-sm text-pretty text-muted-foreground">
          Saved. It replaces your earlier recording of this word.
        </p>
      )}
      {error && <FieldError>{error}</FieldError>}
    </div>
  )
}
