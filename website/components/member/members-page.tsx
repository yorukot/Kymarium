"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { MoreVertical, Shield, Trash } from "lucide-react";

import { useUser } from "@/components/context/user-context";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import { ApiError } from "@/lib/api/client";
import { removeTeamMember, updateTeamMemberRole } from "@/lib/api/team-members";
import { cancelTeamInvite } from "@/lib/api/team-invites";
import { formatRelativeTime } from "@/lib/parsers/datetime";
import { humanizeIdentifier } from "@/lib/parsers/strings";
import type { TeamInvite } from "@/lib/schemas/team-invite";
import type { TeamMember, TeamRole } from "@/lib/schemas/team-member";

import InviteMemberDialog from "@/components/member/invite-member-dialog";

const assignableRoles = ["admin", "member", "viewer"] as const;

function roleRank(role: string): number {
  switch (role) {
    case "owner":
      return 3;
    case "admin":
      return 2;
    case "member":
      return 1;
    case "viewer":
      return 0;
    default:
      return -1;
  }
}

function canInvite(teamRole: string) {
  return teamRole === "owner" || teamRole === "admin";
}

function canRemoveMember(currentRole: TeamRole, targetRole: TeamRole) {
  if (currentRole === "owner") return true;
  if (currentRole !== "admin") return false;

  // Keep this conservative: admins can remove members/viewers, but not admins/owner.
  const targetRank = roleRank(String(targetRole));
  if (targetRank < 0) return false;
  return targetRank <= roleRank("member");
}

function canEditMemberRole(currentRole: TeamRole, targetRole: TeamRole) {
  if (currentRole !== "owner" && currentRole !== "admin") return false;
  if (targetRole === "owner") return false;

  if (currentRole === "admin" && targetRole === "admin") return false;

  return true;
}

function initialsFromName(value: string) {
  const parts = value
    .split(/\s+/)
    .map((p) => p.trim())
    .filter(Boolean);
  const first = parts[0]?.[0] ?? "";
  const second = parts[1]?.[0] ?? "";
  const initial = (first + second).toUpperCase();
  return initial || "?";
}

export default function MembersPageClient({
  teamID,
  teamRole,
  members,
  invites,
}: {
  teamID: string;
  teamRole: TeamRole;
  members: TeamMember[];
  invites: TeamInvite[];
}) {
  const user = useUser();
  const router = useRouter();
  const [busyID, setBusyID] = React.useState<string | null>(null);

  const canSeeInvites = canInvite(String(teamRole));
  const pendingInvites = React.useMemo(
    () => invites.filter((invite) => invite.status === "pending"),
    [invites],
  );

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-between items-center">
        <div className="flex flex-col gap-1">
          <span className="text-xl font-bold">Members</span>
          <span className="text-sm text-muted-foreground">
            Manage who has access to this team
          </span>
        </div>
        {canSeeInvites ? (
          <InviteMemberDialog teamID={teamID} onInvited={() => router.refresh()} />
        ) : null}
      </div>

      <div className="flex items-center justify-between">
        <span className="text-sm font-medium">Members</span>
      </div>

      {members.length === 0 ? (
        <div className="rounded-md border border-dashed px-3 py-2 text-sm text-muted-foreground">
          No members found.
        </div>
      ) : (
        <div className="flex flex-col gap-3">
          {members.map((member) => {
            const isSelf = member.userID === user.id;
            const displayName = member.displayName?.trim() || member.email || "Member";
            const canRemove = !isSelf && canRemoveMember(teamRole, member.role);
            const canEditRole = !isSelf && canEditMemberRole(teamRole, member.role);
            const showActions = canRemove || canEditRole;

            return (
              <Card key={member.id} className="py-1.5 px-3">
                <CardContent className="flex items-center justify-between gap-3 p-0">
                  <div className="flex min-w-0 items-center gap-3">
                    <Avatar className="size-8">
                      {member.avatar ? (
                        <AvatarImage src={member.avatar} alt="" />
                      ) : null}
                      <AvatarFallback>{initialsFromName(displayName)}</AvatarFallback>
                    </Avatar>

                    <div className="flex min-w-0 flex-col gap-1">
                      <div className="flex min-w-0 items-center gap-2">
                        <span className="min-w-0 flex-1 truncate font-medium">
                          {displayName}
                        </span>
                        {isSelf ? <Badge variant="secondary">You</Badge> : null}
                        <Badge variant="outline" className="capitalize">
                          {humanizeIdentifier(String(member.role))}
                        </Badge>
                      </div>
                      {member.email ? (
                        <span className="text-sm text-muted-foreground truncate">
                          {member.email}
                        </span>
                      ) : null}
                    </div>
                  </div>

                  <div className="flex flex-none items-center gap-3">
                    {showActions ? (
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button
                            size="icon"
                            variant="ghost"
                            aria-label="Actions"
                            disabled={busyID === member.id}
                          >
                            <MoreVertical size={16} />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="start">
                          <DropdownMenuGroup>
                            {canEditRole ? (
                              <DropdownMenuSub>
                                <DropdownMenuSubTrigger>
                                  <Shield />
                                  Change role
                                </DropdownMenuSubTrigger>
                                <DropdownMenuSubContent>
                                  <DropdownMenuRadioGroup
                                    value={String(member.role)}
                                    onValueChange={async (value) => {
                                      if (value === String(member.role)) return;
                                      setBusyID(member.id);
                                      try {
                                        await updateTeamMemberRole(teamID, member.userID, value);
                                        toast.success("Member role updated");
                                        router.refresh();
                                      } catch (error) {
                                        if (error instanceof ApiError) {
                                          toast.error(error.message);
                                        } else {
                                          toast.error("Failed to update member role");
                                        }
                                      } finally {
                                        setBusyID(null);
                                      }
                                    }}
                                  >
                                    {assignableRoles.map((role) => (
                                      <DropdownMenuRadioItem key={role} value={role}>
                                        {humanizeIdentifier(role)}
                                      </DropdownMenuRadioItem>
                                    ))}
                                  </DropdownMenuRadioGroup>
                                </DropdownMenuSubContent>
                              </DropdownMenuSub>
                            ) : null}

                            {canEditRole && canRemove ? <DropdownMenuSeparator /> : null}

                            <DropdownMenuItem
                              variant="destructive"
                              onSelect={async () => {
                                if (!canRemove) return;
                                const confirmed = window.confirm(
                                  `Remove "${displayName}" from the team?`,
                                );
                                if (!confirmed) return;

                                setBusyID(member.id);
                                try {
                                  await removeTeamMember(teamID, member.userID);
                                  toast.success("Member removed");
                                  router.refresh();
                                } catch (error) {
                                  if (error instanceof ApiError) {
                                    toast.error(error.message);
                                  } else {
                                    toast.error("Failed to remove member");
                                  }
                                } finally {
                                  setBusyID(null);
                                }
                              }}
                              disabled={!canRemove}
                            >
                              <Trash />
                              Remove
                            </DropdownMenuItem>
                          </DropdownMenuGroup>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    ) : null}
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}

      {canSeeInvites ? (
        <div className="flex flex-col gap-3">
          <div className="flex items-center justify-between">
            <span className="text-sm font-medium">Pending invites</span>
          </div>

          {pendingInvites.length === 0 ? (
            <div className="rounded-md border border-dashed px-3 py-2 text-sm text-muted-foreground">
              No pending invites.
            </div>
          ) : (
            <div className="flex flex-col gap-3">
              {pendingInvites.map((invite) => (
                <Card key={invite.id} className="py-1.5 px-3">
                  <CardContent className="flex items-center justify-between gap-3 p-0">
                    <div className="flex min-w-0 flex-col gap-1">
                      <div className="flex min-w-0 items-center gap-2">
                        <span className="min-w-0 flex-1 truncate font-medium">
                          {invite.invitedEmail}
                        </span>
                        <Badge variant="outline" className="capitalize">
                          {humanizeIdentifier(String(invite.role))}
                        </Badge>
                      </div>
                      <span className="text-sm text-muted-foreground">
                        Expires {formatRelativeTime(invite.expiresAt)}
                      </span>
                    </div>

                    <div className="flex flex-none items-center gap-3">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button
                            size="icon"
                            variant="ghost"
                            aria-label="Actions"
                            disabled={busyID === invite.id}
                          >
                            <MoreVertical size={16} />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="start">
                          <DropdownMenuGroup>
                            <DropdownMenuItem
                              variant="destructive"
                              onSelect={async () => {
                                const confirmed = window.confirm(
                                  `Cancel invite for "${invite.invitedEmail}"?`,
                                );
                                if (!confirmed) return;

                                setBusyID(invite.id);
                                try {
                                  await cancelTeamInvite(teamID, invite.id);
                                  toast.success("Invite canceled");
                                  router.refresh();
                                } catch (error) {
                                  if (error instanceof ApiError) {
                                    toast.error(error.message);
                                  } else {
                                    toast.error("Failed to cancel invite");
                                  }
                                } finally {
                                  setBusyID(null);
                                }
                              }}
                            >
                              <Trash />
                              Cancel invite
                            </DropdownMenuItem>
                          </DropdownMenuGroup>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>
      ) : (
        <div className="text-sm text-muted-foreground">
          You don&apos;t have permission to invite or manage team members.
        </div>
      )}
    </div>
  );
}
