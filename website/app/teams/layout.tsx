import type { ReactNode } from "react"
import { redirect } from "next/navigation"

import {
  TeamsProvider,
  type TeamSummary,
} from "@/components/teams/teams-context"
import { buildCookieHeader } from "@/lib/api/cookies"

export const dynamic = "force-dynamic"

type TeamsResponse = {
  message?: string
  data?: TeamSummary[]
}

function debugLog(message: string, details?: Record<string, unknown>) {
  if (process.env.NODE_ENV === "production") return
  if (details) {
    console.info(`[teams/layout] ${message}`, details)
  } else {
    console.info(`[teams/layout] ${message}`)
  }
}

async function fetchTeams(): Promise<TeamSummary[]> {
  const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_API_BASE_URL")
  }

  const cookieHeader = await buildCookieHeader()
  debugLog("request cookies", {
    hasCookieHeader: Boolean(cookieHeader),
    cookieHeaderLength: cookieHeader.length,
  })

  const res = await fetch(`${apiBase}/api/teams`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  })
  debugLog("teams response", { status: res.status })

  if (res.status === 401) {
    redirect("/login?next=/teams")
  }

  if (!res.ok) {
    throw new Error("Failed to load teams")
  }

  const body = (await res.json()) as TeamsResponse
  if (!Array.isArray(body?.data)) {
    return []
  }

  return body.data.map((team) => ({
    id: String(team.id),
    name: team.name,
    role: team.role,
  }))
}

export default async function TeamsLayout({
  children,
}: {
  children: ReactNode
}) {
  const teams = await fetchTeams()
  debugLog("teams loaded", { count: teams.length })

  return <TeamsProvider teams={teams}>{children}</TeamsProvider>
}
