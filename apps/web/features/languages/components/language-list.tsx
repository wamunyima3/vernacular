import { Badge } from "@/components/ui/badge"
import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

import { getLanguages } from "../get-languages"
import { groupByWave, type Language } from "../group-by-wave"
import { LanguageGrid } from "./language-grid"

/** Languages grouped by rollout wave, fetched from GET /v1/languages at request time. */
export async function LanguageList() {
  const result = await getLanguages()

  if (!result.ok) {
    return (
      <Card role="alert">
        <CardHeader>
          <CardTitle>Languages could not be loaded</CardTitle>
          <CardDescription className="text-pretty">
            {result.message} Start the API with <code className="font-mono">make dev-api</code> and
            reload the page.
          </CardDescription>
        </CardHeader>
      </Card>
    )
  }

  return (
    <div className="flex flex-col gap-8">
      {groupByWave(result.languages).map(({ wave, languages }) => (
        <div key={wave} className="flex flex-col gap-3">
          <h3 className="text-sm font-medium text-muted-foreground">
            Wave <span className="tabular-nums">{wave}</span>
          </h3>
          <LanguageGrid>
            {languages.map((language) => (
              <LanguageCard key={language.code} language={language} />
            ))}
          </LanguageGrid>
        </div>
      ))}
    </div>
  )
}

function LanguageCard({ language }: { language: Language }) {
  const active = language.status === "active"
  return (
    <Card size="sm">
      <CardHeader>
        <CardTitle className="truncate">{language.name}</CardTitle>
        <CardDescription className="truncate">
          {language.own_name} · <span className="font-mono">{language.code}</span>
        </CardDescription>
        <Badge variant={active ? "default" : "outline"} className="mt-1">
          {active ? "Active" : "Planned"}
        </Badge>
      </CardHeader>
    </Card>
  )
}
