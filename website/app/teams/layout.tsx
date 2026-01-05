import type { ReactNode } from "react"
import { redirect } from "next/navigation"

import {
  TeamsProvider,
  type TeamSummary,
} from "@/components/context/teams-context"
import {
  UserProvider,
  type UserSummary,
} from "@/components/context/user-context"
import { buildCookieHeader } from "@/lib/api/cookies"

export const dynamic = "force-dynamic"

type TeamsResponse = {
  message?: string
  data?: TeamSummary[]
}

type UserResponse = {
  message?: string
  data?: {
    id: string
    display_name: string
    avatar?: string | null
  }
}

type AccountSummary = {
  email: string
  is_primary: boolean
}

type AccountsResponse = {
  message?: string
  data?: AccountSummary[]
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

async function fetchUser(): Promise<UserSummary> {
  const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_API_BASE_URL")
  }

  const cookieHeader = await buildCookieHeader()
  debugLog("user request cookies", {
    hasCookieHeader: Boolean(cookieHeader),
    cookieHeaderLength: cookieHeader.length,
  })

  const [userRes, accountsRes] = await Promise.all([
    fetch(`${apiBase}/api/users/me`, {
      method: "GET",
      headers: cookieHeader ? { cookie: cookieHeader } : undefined,
      cache: "no-store",
    }),
    fetch(`${apiBase}/api/users/me/account`, {
      method: "GET",
      headers: cookieHeader ? { cookie: cookieHeader } : undefined,
      cache: "no-store",
    }),
  ])

  debugLog("user response", {
    status: userRes.status,
    accountStatus: accountsRes.status,
  })

  if (userRes.status === 401 || accountsRes.status === 401) {
    redirect("/login?next=/teams")
  }

  if (!userRes.ok) {
    throw new Error("Failed to load user")
  }

  const userBody = (await userRes.json()) as UserResponse
  const userData = userBody?.data
  if (!userData) {
    throw new Error("Missing user data")
  }

  let email = ""
  if (accountsRes.ok) {
    const accountsBody = (await accountsRes.json()) as AccountsResponse
    const accounts = Array.isArray(accountsBody?.data)
      ? accountsBody.data
      : []
    const primary = accounts.find((account) => account.is_primary)
    email = primary?.email ?? accounts[0]?.email ?? ""
  }

  return {
    id: String(userData.id),
    displayName: userData.display_name ?? "",
    email,
    avatar: userData.avatar ?? null,
  }
}

export default async function TeamsLayout({
  children,
}: {
  children: ReactNode
}) {
  const [teams, user] = await Promise.all([fetchTeams(), fetchUser()])
  debugLog("teams loaded", { count: teams.length })

  return (
    <UserProvider user={user}>
      <TeamsProvider teams={teams}>{children}</TeamsProvider>
    </UserProvider>
  )
}
