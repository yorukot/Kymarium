import type { Metadata } from "next"
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

export const metadata: Metadata = {
  title: {
    default: "Teams",
    template: "%s | Kymarium",
  },
}

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

async function fetchTeams(): Promise<TeamSummary[]> {
  const apiBase = process.env.NEXT_PUBLIC_BACKEND_URL
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_BACKEND_URL")
  }

  const cookieHeader = await buildCookieHeader()
  const res = await fetch(`${apiBase}/api/teams`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  })

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
  const apiBase = process.env.NEXT_PUBLIC_BACKEND_URL
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_BACKEND_URL")
  }

  const cookieHeader = await buildCookieHeader()

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

  return (
    <UserProvider user={user}>
      <TeamsProvider teams={teams}>{children}</TeamsProvider>
    </UserProvider>
  )
}
