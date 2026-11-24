import { ConfirmAction } from "@/components/confirm-action"
import { formatDate } from "@/lib/format"

import { apiForSession } from "../../auth/session"
import { withdrawConsent } from "../actions"

/**
 * The signed-in user's consent to be recorded, with a way to withdraw it (NFR-6, US-33).
 * Shows nothing for users without a consent record or while audio is switched off.
 */
export async function SpeakerConsent() {
  const api = await apiForSession("/account")
  const { data } = await api.GET("/v1/me/speaker")
  if (!data) return null
  return (
    <section className="flex flex-col gap-3" aria-labelledby="recordings">
      <h2 id="recordings" className="text-xl font-semibold text-balance">
        Recordings
      </h2>
      {data.withdrawn_at ? (
        <p className="text-sm text-pretty text-muted-foreground">
          You withdrew your consent on {formatDate(data.withdrawn_at)} and your recordings were
          deleted. To record again, sign a new consent form with an admin.
        </p>
      ) : (
        <>
          <p className="text-sm text-pretty text-muted-foreground">
            You agreed on {formatDate(data.consented_at)} to record pronunciations in{" "}
            {data.languages.join(", ")}. You can withdraw at any time: your recordings are then
            deleted.
          </p>
          <ConfirmAction
            className="self-start"
            label="Withdraw consent"
            title="Withdraw your consent?"
            description="All your recordings are deleted and you can no longer record. An admin can record a new consent form if you change your mind."
            action={withdrawConsent}
            destructive
          />
        </>
      )}
    </section>
  )
}
