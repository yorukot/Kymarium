"use client";

import {
  BadgeCheck,
  Bell,
  ChevronsUpDown,
  CreditCard,
  Mail,
  LogOut,
  Sparkles,
} from "lucide-react";
import { createAvatar } from "@dicebear/core";
import { thumbs } from "@dicebear/collection";
import * as React from "react";
import { toast } from "sonner";
import { useParams, useRouter } from "next/navigation";

import { Avatar, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar";
import { Spinner } from "@/components/ui/spinner";
import { useUser } from "@/components/context/user-context";
import { logout } from "@/lib/api/auth";
import { ApiError } from "@/lib/api/client";
import { listMyPendingInvites, respondToTeamInvite } from "@/lib/api/user-invites";
import { formatRelativeTime } from "@/lib/parsers/datetime";
import { humanizeIdentifier } from "@/lib/parsers/strings";
import type { TeamInviteWithTeam } from "@/lib/schemas/team-invite";

export function NavUser() {
  const user = useUser();
  const { isMobile } = useSidebar();
  const router = useRouter();
  const params = useParams<{ teamID: string }>();
  const teamID = params.teamID;

  const [inviteDialogOpen, setInviteDialogOpen] = React.useState(false);
  const [invites, setInvites] = React.useState<TeamInviteWithTeam[]>([]);
  const [invitesLoading, setInvitesLoading] = React.useState(false);
  const [busyInviteID, setBusyInviteID] = React.useState<string | null>(null);
  const [loggingOut, setLoggingOut] = React.useState(false);

  const avatarSeed = user.id;
  const avatarSrc = React.useMemo(
    () =>
      createAvatar(thumbs, {
        seed: avatarSeed,
        size: 64,
      }).toDataUri(),
    [avatarSeed],
  );

  const pendingInviteCount = invites.length;

  const loadInvites = React.useCallback(async () => {
    setInvitesLoading(true);
    try {
      const data = await listMyPendingInvites();
      setInvites(data);
    } catch (error) {
      if (error instanceof ApiError) {
        toast.error(error.message);
        return;
      }
      toast.error("Failed to fetch invites");
    } finally {
      setInvitesLoading(false);
    }
  }, []);

  React.useEffect(() => {
    void loadInvites();
  }, [loadInvites]);

  React.useEffect(() => {
    if (!inviteDialogOpen) return;
    void loadInvites();
  }, [inviteDialogOpen, loadInvites]);

  const respond = React.useCallback(
    async (invite: TeamInviteWithTeam, status: "accepted" | "rejected") => {
      setBusyInviteID(invite.id);
      const toastId = toast.loading(
        status === "accepted" ? "Accepting invite..." : "Rejecting invite...",
      );
      try {
        await respondToTeamInvite(invite.teamID, invite.id, status);
        toast.success(
          status === "accepted" ? "Invite accepted" : "Invite rejected",
          { id: toastId },
        );
        setInvites((prev) => prev.filter((x) => x.id !== invite.id));
        router.refresh();
      } catch (error) {
        if (error instanceof ApiError) {
          toast.error(error.message, { id: toastId });
          return;
        }
        toast.error("Failed to update invite", { id: toastId });
      } finally {
        setBusyInviteID(null);
      }
    },
    [router],
  );

  const handleLogout = React.useCallback(async () => {
    if (loggingOut) return;

    setLoggingOut(true);
    const toastId = toast.loading("Logging out...");

    try {
      await logout();
      toast.success("Logged out", { id: toastId });

      const currentPath =
        window.location.pathname +
        window.location.search +
        window.location.hash;

      router.replace(`/login?next=${encodeURIComponent(currentPath)}`);
    } catch (error) {
      if (error instanceof ApiError) {
        toast.error(error.message, { id: toastId });
      } else {
        toast.error("Logout failed", { id: toastId });
      }
      setLoggingOut(false);
    }
  }, [loggingOut, router]);

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              id="nav-user-trigger"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <Avatar className="h-8 w-8 rounded-lg">
                <AvatarImage
                  src={user.avatar || avatarSrc}
                  alt={user.displayName}
                />
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{user.displayName}</span>
                <span className="truncate text-xs">{user.email}</span>
              </div>
              <ChevronsUpDown className="ml-auto size-4" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
            side={isMobile ? "bottom" : "right"}
            align="end"
            sideOffset={4}
          >
            <DropdownMenuLabel className="p-0 font-normal">
              <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                <Avatar className="h-8 w-8 rounded-lg">
                  <AvatarImage src={avatarSrc} alt={user.displayName} />
                </Avatar>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">
                    {user.displayName}
                  </span>
                  <span className="truncate text-xs">{user.email}</span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <Sparkles />
                Upgrade to Pro
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem
                onSelect={() => setInviteDialogOpen(true)}
                className="relative pr-10"
              >
                <Mail />
                Pending invites
                {pendingInviteCount > 0 ? (
                  <Badge
                    variant="secondary"
                    className="absolute top-1.5 right-2 h-5 min-w-5 px-1"
                  >
                    {pendingInviteCount}
                  </Badge>
                ) : null}
              </DropdownMenuItem>
              <DropdownMenuItem
                disabled={!teamID}
                onSelect={() => router.push(`/teams/${teamID}/settings/account`)}
              >
                <BadgeCheck />
                Account
              </DropdownMenuItem>
              <DropdownMenuItem>
                <CreditCard />
                Billing
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Bell />
                Notifications
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              disabled={loggingOut}
              onSelect={() => void handleLogout()}
            >
              {loggingOut ? <Spinner /> : <LogOut />}
              {loggingOut ? "Logging out..." : "Log out"}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        <Dialog open={inviteDialogOpen} onOpenChange={setInviteDialogOpen}>
          <DialogContent className="sm:max-w-xl">
            <DialogHeader>
              <DialogTitle>Pending invites</DialogTitle>
              <DialogDescription>
                Invites for your account. Accept to join the team or reject to ignore it.
              </DialogDescription>
            </DialogHeader>

            <div className="max-h-[60vh] space-y-2 overflow-auto">
              {invitesLoading ? (
                <div className="text-muted-foreground flex items-center justify-center gap-2 py-8 text-sm">
                  <Spinner />
                  Loading invites...
                </div>
              ) : invites.length === 0 ? (
                <div className="text-muted-foreground py-8 text-center text-sm">
                  No pending invites.
                </div>
              ) : (
                invites.map((invite) => (
                  <Card key={invite.id} className="gap-3 py-4 px-4">
                    <div className="flex items-start justify-between gap-3">
                      <div className="min-w-0">
                        <div className="truncate text-sm font-medium">
                          {invite.teamName}
                        </div>
                        <div className="text-muted-foreground mt-1 text-xs">
                          Role: {humanizeIdentifier(String(invite.role))} • Expires{" "}
                          {formatRelativeTime(invite.expiresAt)}
                        </div>
                      </div>

                      <div className="flex shrink-0 gap-2">
                        <Button
                          size="sm"
                          variant="outline"
                          disabled={busyInviteID === invite.id}
                          onClick={() => respond(invite, "rejected")}
                        >
                          {busyInviteID === invite.id ? <Spinner /> : null}
                          Reject
                        </Button>
                        <Button
                          size="sm"
                          disabled={busyInviteID === invite.id}
                          onClick={() => respond(invite, "accepted")}
                        >
                          {busyInviteID === invite.id ? <Spinner /> : null}
                          Accept
                        </Button>
                      </div>
                    </div>
                  </Card>
                ))
              )}
            </div>

            <DialogFooter>
              <Button variant="outline" onClick={() => setInviteDialogOpen(false)}>
                Close
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
