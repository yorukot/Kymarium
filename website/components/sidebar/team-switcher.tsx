"use client"

import * as React from "react"
import { ChevronsUpDown, Plus } from "lucide-react"
import { useParams, useRouter } from "next/navigation"
import Image from "next/image"
import { createAvatar } from "@dicebear/core"
import { shapes } from "@dicebear/collection"

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar"
import { useTeams } from "@/components/context/teams-context"

function formatRole(role: string) {
  if (!role) return ""
  return role
    .split(/[-_]/)
    .filter(Boolean)
    .map((part) => part[0]?.toUpperCase() + part.slice(1))
    .join(" ")
}

function avatarFor(seed: string, size: number) {
  return createAvatar(shapes, {
    seed,
    size,
  }).toDataUri()
}

export function TeamSwitcher() {
  const teams = useTeams()
  const { isMobile } = useSidebar()
  const router = useRouter()
  const params = useParams<{ teamID?: string }>()
  const teamID = params.teamID
  const activeTeam = React.useMemo(() => {
    if (teamID) {
      const found = teams.find((team) => team.id === teamID)
      if (found) return found
    }
    return teams[0]
  }, [teams, teamID])

  if (!activeTeam) {
    return null
  }

  const activeRole = formatRole(activeTeam.role)

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
                <Image
                  src={avatarFor(activeTeam.id, 32)}
                  alt=""
                  width={32}
                  height={32}
                  unoptimized
                  className="rounded-lg"
                />
              </div>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{activeTeam.name}</span>
                <span className="truncate text-xs">{activeRole}</span>
              </div>
              <ChevronsUpDown className="ml-auto" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
            align="start"
            side={isMobile ? "bottom" : "right"}
            sideOffset={4}
          >
            <DropdownMenuLabel className="text-muted-foreground text-xs">
              Teams
            </DropdownMenuLabel>
            {teams.map((team) => (
              <DropdownMenuItem
                key={team.id}
                onSelect={() => router.push(`/teams/${team.id}/monitors`)}
                className="gap-2 p-2"
              >
                <div className="flex size-6 items-center justify-center rounded-md border">
                  <Image
                    src={avatarFor(team.id, 24)}
                    alt=""
                    width={24}
                    height={24}
                    unoptimized
                    className="rounded"
                  />
                </div>
                {team.name}
              </DropdownMenuItem>
            ))}
            <DropdownMenuSeparator />
            <DropdownMenuItem
              className="gap-2 p-2"
              onSelect={() => router.push("/teams/new-team")}
            >
              <div className="flex size-6 items-center justify-center rounded-md border bg-transparent">
                <Plus className="size-4" />
              </div>
              <div className="text-muted-foreground font-medium">Add team</div>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
